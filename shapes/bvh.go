// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	"sort"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
)

// BVH object information
type BvhObjectInfo struct {
	idx      int   // Index of object in the object array
	bounds   Box   // Object bounds in world space
	centroid Tuple // Centroid of object bounds
}

// BVH tree node used in the build phase
type BvhBuildNode struct {
	bounds    Box // Bounds for all objects inside the node
	child     [2]*BvhBuildNode
	objIdx    int // Index of object in the object array
	objCount  int // How many objects this node contains, consecutive starting from index
	splitAxis int // Which axis was used to partition the node
}

// BVH node used in the optimized linear (array) representation
type BvhLinearNode struct {
	bounds   Box
	index    int // Index of object if leaf, index of second child if interior
	objCount int // How many objects in this node (0 for interior)
	axis     int
}

type BvhSplitMethod int

const (
	BvhSplitNone        BvhSplitMethod = iota // List is just split without applying any geometric criteria
	BvhSplitMiddle                            // Elements above and below the midpoint of the splitting axis
	BvhSplitEqualCounts                       // Ordered by centroid along the splitting axis, then split in two equally sized sets
	BvhSplitSAH                               // A Surface Area Heuristic is used to estimate and minimize the later cost of intersections
)

type BvhBucketInfo struct {
	bounds Box     // Bounds for all objects in bucket
	count  int     // How many objects in this bucket
	cost   float64 // Cost of splitting _after_ this bucket
}

const BvhMaxObjectsPerNode = 16

const BvhNumBuckets = 8

const BvhBboxIntersectionCost = 0.3 // Relative to the cost of intersecting a primitive shape, which is set to 1 as a reference

// BuildBVH build a Bounding Volume Hierarchy for the group,
// using the algorithms described in Physically Based Rendering
func (g *Group) BuildBVH() {
	objInfo := make([]BvhObjectInfo, len(g.members))
	orderedObjects := make([]Groupable, 0, len(g.members))

	// Phase 1: collect info about all objects and build bounds
	for i, s := range g.members {
		bbox := s.Bounds(). // Bounds in object local space
					Transform(s.Transform()). // Bounds in group local space
					Transform(g.Transform())  // Bounds in world space

		objInfo[i] = BvhObjectInfo{
			i,
			bbox,
			Point(
				(bbox.Min.X+bbox.Max.X)/2,
				(bbox.Min.Y+bbox.Max.Y)/2,
				(bbox.Min.Z+bbox.Max.Z)/2,
			),
		}
	}

	// Phase 2: build tree
	method := BvhSplitSAH

	totalNodes := 0

	var bvhRecursiveBuild func(int, int, int) *BvhBuildNode

	bvhRecursiveBuild = func(start, end, level int) *BvhBuildNode {
		node := BvhBuildNode{}
		totalNodes++

		if start < end {
			objInfo[start], objInfo[end-1] = objInfo[end-1], objInfo[start]
		}

		// Compute bound of all objects inside this node
		bounds := NewBox(PointAtInfinity(+1), PointAtInfinity(-1))
		for i := start; i < end; i++ {
			bounds = bounds.Union(objInfo[i].bounds)
		}

		numObjects := end - start

		if numObjects > 1 {
			// Inner node, get the centroids bounds
			centroidBounds := NewBox(PointAtInfinity(+1), PointAtInfinity(-1))
			for i := start; i < end; i++ {
				centroid := objInfo[i].centroid
				centroidBounds = centroidBounds.Union(NewBox(centroid, centroid))
			}

			// Get the centroid maximum extend (longest axis)
			diagonal := centroidBounds.Max.Sub(centroidBounds.Min)
			dim := 2 // Z-axis
			if diagonal.X > diagonal.Y && diagonal.X > diagonal.Z {
				dim = 0 // X-axis
			} else if diagonal.Y > diagonal.Z {
				dim = 1 // Y-axis
			}

			// Partition
			mid := (start + end) / 2

			// Debugln("INTERIOR start=",start, ",end=", end, "mid=", mid, ", dim=", dim, "d=", diagonal)

			if !FloatEqual(centroidBounds.Max.CompByIdx(dim), centroidBounds.Min.CompByIdx(dim)) {
				// Choose dim according to split method
				switch method {
				case BvhSplitMiddle:
					// Partition objects based on whether they are below or above the midpoint of the splitting axis
					pmid := (centroidBounds.Min.CompByIdx(dim) + centroidBounds.Max.CompByIdx(dim)) / 2

					e := end - 1
					for s := start; s < e; {
						if objInfo[s].centroid.CompByIdx(dim) < pmid {
							// Move element to the end
							objInfo[e], objInfo[s] = objInfo[s], objInfo[e]
							e--
						} else {
							// Already in the proper side of the partition
							s++
						}
					}

					mid = e + 1

					if mid != start && mid != end {
						break
					}

					// As we failed to partition the box, let's try another method
					fallthrough
				case BvhSplitEqualCounts:
					// Partition objects into equally sized subsets
					// Note: we don't really need to full sort here and we could achieve our partition in O(n) time,
					// but the implementation is not easy so let's have it this way for now
					sort.Slice(objInfo[start:end], func(i, j int) bool {
						return objInfo[i].centroid.CompByIdx(dim) < objInfo[j].centroid.CompByIdx(dim)
					})

					mid = (start + end) / 2
				case BvhSplitSAH:
					if numObjects > 2 {
						// Allocate buckets
						buckets := [BvhNumBuckets]BvhBucketInfo{}

						for i := range buckets {
							buckets[i].bounds = NewBox(PointAtInfinity(+1), PointAtInfinity(-1))
						}

						// Initialize bucket info
						for i := start; i < end; i++ {
							b := int(BvhNumBuckets * centroidBounds.Offset(objInfo[i].centroid).CompByIdx(dim))

							if b == BvhNumBuckets {
								b--
							}

							// Note: if we get an "index out of bounds" panic here, then probably b is NaN
							// because the group contains planes, which have infinite bounds
							buckets[b].count++
							buckets[b].bounds = buckets[b].bounds.Union(objInfo[i].bounds)
						}

						// Compute splitting cost
						for i := 0; i < BvhNumBuckets-1; i++ {
							b0 := NewBox(PointAtInfinity(+1), PointAtInfinity(-1))
							b1 := NewBox(PointAtInfinity(+1), PointAtInfinity(-1))
							count0 := 0
							count1 := 0

							for j := 0; j <= i; j++ {
								b0 = b0.Union(buckets[j].bounds)
								count0 += buckets[j].count
							}

							for j := i + 1; j < BvhNumBuckets; j++ {
								b1 = b1.Union(buckets[j].bounds)
								count1 += buckets[j].count
							}

							buckets[i].cost = BvhBboxIntersectionCost + (float64(count0)*b0.SurfaceArea()+float64(count1)*b1.SurfaceArea())/bounds.SurfaceArea()
						}

						// Find bucket to split at that minimizes cost
						minCost := buckets[0].cost
						minBucket := 0

						for i := 1; i < BvhNumBuckets-1; i++ {
							if buckets[i].cost < minCost {
								minCost = buckets[i].cost
								minBucket = i
							}
						}

						// Split if convenient, otherwise create a leaf node
						leafCost := float64(numObjects) // Cost of intersecting an object (primitive) is set to 1 (see also BvhBboxIntersectionCost)

						if numObjects > BvhMaxObjectsPerNode || minCost < leafCost {
							// Split
							e := end - 1
							for s := start; s < e; {
								b := int(BvhNumBuckets * centroidBounds.Offset(objInfo[s].centroid).CompByIdx(dim))

								if b == BvhNumBuckets {
									b--
								}

								if b > minBucket {
									// Move element to the end
									objInfo[e], objInfo[s] = objInfo[s], objInfo[e]
									e--
								} else {
									// Already in the proper side of the partition
									s++
								}
							}

							mid = e + 1 // Will create a leaf if no partition is found
						} else {
							// Create a leaf node
							mid = end
						}
					}
					// ...else with just two objects we just split in the middle
				}

				// Build an interior node
				if mid != end {
					node.splitAxis = dim
					node.child[0] = bvhRecursiveBuild(start, mid, level+1)
					node.child[1] = bvhRecursiveBuild(mid, end, level+1)
					node.bounds = node.child[0].bounds.Union(node.child[1].bounds)

					return &node
				}
				// ...else we don't have a split point so make it a leaf
			}
			// ...else all centroid points are at the same position (rare!) and we can't handle this case, make it a leaf
		}
		// ...else there is only one object, so make it a leaf

		// Leaf
		// Debugln("LEAF")
		node.bounds = bounds
		node.objIdx = len(orderedObjects)
		node.objCount = numObjects
		for i := start; i < end; i++ {
			orderedObjects = append(orderedObjects, g.members[objInfo[i].idx])
		}

		return &node
	}

	root := bvhRecursiveBuild(0, len(objInfo), 0)

	g.members = orderedObjects

	// Phase 3: flatten tree
	nodes := make([]BvhLinearNode, totalNodes)
	offset := 0

	var flattenBvhTree func(*BvhBuildNode) int

	flattenBvhTree = func(node *BvhBuildNode) int {
		off := offset
		offset++

		nodes[off].bounds = node.bounds

		if node.objCount > 0 {
			nodes[off].index = node.objIdx
			nodes[off].objCount = node.objCount
		} else {
			nodes[off].axis = node.splitAxis
			flattenBvhTree(node.child[0])
			nodes[off].index = flattenBvhTree(node.child[1])
		}

		return off
	}

	flattenBvhTree(root)

	g.bvhNodes = nodes
}

// AddIntersectionsBvh checks for intersections between a ray and all objects
// in the group, using a BVH for performance
func (g *Group) AddIntersectionsBvh(ray Ray, xs *Intersections) {
	rayInObjectSpace := ray.Transform(g.Tinverse)

	toVisitOffset := 0
	currentNodeIndex := 0
	nodesToVisit := [64]int{}

	ray.Direction.X = 1 / ray.Direction.X // Precompute inverse direction
	ray.Direction.Y = 1 / ray.Direction.Y
	ray.Direction.Z = 1 / ray.Direction.Z

	for {
		node := &(g.bvhNodes[currentNodeIndex])

		if node.bounds.IntersectsInvDir(ray) {
			if node.objCount > 0 {
				// Leaf: test all objects
				for i := 0; i < node.objCount; i++ {
					s := g.members[node.index+i]
					s.AddIntersections(rayInObjectSpace, xs)
				}

				if toVisitOffset == 0 {
					break
				}

				toVisitOffset--
				currentNodeIndex = nodesToVisit[toVisitOffset]
			} else {
				// Interiori: put one child on stack and advance to the other
				currentNodeIndex = currentNodeIndex + 1
				nodesToVisit[toVisitOffset] = node.index // Index of second child
				toVisitOffset++
			}
		} else {
			if toVisitOffset == 0 {
				break
			}

			toVisitOffset--
			currentNodeIndex = nodesToVisit[toVisitOffset]
		}
	}
}
