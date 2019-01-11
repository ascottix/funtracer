// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestSbtBox(t *testing.T) {
	TestWithImage(t)

	scene := `
SBT-raytracer 1.0

// Slightly adapted from box.ray (removed texture map)

camera {
	position = (4,0,0);
	viewdir = (-1,0,0);
	aspectratio = 1;
	updir = (0,1,0);
}

point_light {
	position = (4, 4, 0);
	color = (.5, .5, .5);
	constant_attenuation_coeff= 0.25;
	linear_attenuation_coeff = 0.003372407;
	quadratic_attenuation_coeff = 0.000045492;	
}

directional_light {
	direction = (0, -1, 0);
	colour = (1.0, 1.0, 1.0);
}

directional_light {
	direction = (0,1,0);
	colour = (0.2,0.2,0.2);
}

rotate( 0, 1, 1, -2,
box { 
	material = { 
		diffuse = (0.7, 0, 1.0);
		specular = (0.9,0.4,0.0);
		shininess = 76.8;
	};
})
    `
	s, _ := ParseSbtSceneFromString(scene)

	if len(s.World.Objects) != 1 {
		t.Errorf("num of world object mismatch, found: %d", len(s.World.Objects))
	}

	s.Camera.SetViewSize(400, 400)
	s.Camera.SetFieldOfView(0.5)

	s.World.RenderToPNG(s.Camera, "test_sbt_box.png")
}

func TestSbtEasy2(t *testing.T) {
	TestWithImage(t)

	scene := `
SBT-raytracer 1.0

camera {
  position=( 5.26211,3.09263,-1.38722 );
  viewdir=( -0.853042,-0.383547,0.353853 );
  updir=( 0.345784,0.0924084,0.933753 );
  fov=45;
}
point_light {
  position=( 5.21898,4.68627,-2.04877 );
  color=( 1,1,1 );
  constant_attenuation_coeff= 0.25;
  linear_attenuation_coeff = 0.003372407;
  quadratic_attenuation_coeff = 0.000045492;	
}
translate( 1.02616,-0.365439,-1.13081,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.288315,0.288315,0.288315,
sphere {
  name="";
  material={
    diffuse=( 0.2,0.6,0.75);
    ambient=( 0.2,0.2,0.2);
    specular=( 0,0,0);
    emissive=( 0,0,0);
    shininess=25.6;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( 7.4119e-08,2.78629,-0.116676,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.310878,0.310878,0.310878,
sphere {
  name="";
  material={
    diffuse=( 0.8,0.2,0.5);
    ambient=( 0.2,0.2,0.2);
    specular=( 0,0,0);
    emissive=( 0,0,0);
    shininess=25.6;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( 1.39952,0.365459,2.25247,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.312566,0.312566,0.312566,
sphere {
  name="";
  material={
    diffuse=( 0.5,0.75,0.2);
    ambient=( 0.2,0.2,0.2);
    specular=( 0,0,0);
    emissive=( 0,0,0);
    shininess=25.6;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( -3.21911e-07,2.01869,0.108636,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.670496,0.670496,0.670496,
sphere {
  name="";
  material={
    diffuse=( 0.34,0.07,0.56);
    ambient=( 0.2,0.2,0.2);
    specular=( 0.4,0.4,0.4);
    emissive=( 0,0,0);
    shininess=122.074112;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( 0.990468,0.489687,1.61544,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.636615,0.636615,0.636615,
sphere {
  name="";
  material={
    diffuse=( 0.56,0.24,0.12);
    ambient=( 0.2,0.2,0.2);
    specular=( 0.4,0.4,0.4);
    emissive=( 0,0,0);
    shininess=120.888832;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( 0.709821,-0.124101,-0.637882,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 0.597582,0.597582,0.597582,
sphere {
  name="";
  material={
    diffuse=( 0.04,0.56,0.28);
    ambient=( 0.2,0.2,0.2);
    specular=( 0.4,0.4,0.4);
    emissive=( 0,0,0);
    shininess=124.444416;
    transmissive=( 0,0,0 );
  };
}
 )))
translate( 0.0579858,0.496598,0.550744,
transform( (1,0,0,0), (0,1,0,0), (0,0,1,0), (0,0,0,1),
scale( 1.18902,1.18902,1.18902,
sphere {
  name="";
  material={
    diffuse=( 0.6,0.6,0.6);
    ambient=( 0.2,0.2,0.2);
    specular=( 0.5,0.5,0.5);
    emissive=( 0,0,0);
    shininess=118.518528;
    transmissive=( 0,0,0 );
  };
}
 )))
    `
	s, _ := ParseSbtSceneFromString(scene)

	s.Camera.SetViewSize(400, 400)

	if len(s.World.Objects) != 7 {
		t.Errorf("num of world object mismatch, found: %d", len(s.World.Objects))
	}

	s.World.RenderToPNG(s.Camera, "test_sbt_easy2.png")
}

func TestSbtCsg(t *testing.T) {
	TestWithImage(t)

	data := `
FUN-raytracer 0.1

camera {
  position = ( 0, 1, -8 );
  target = (0, 0, 0);
  updir = ( 0, 1, 0 );
  fovrad = 1;
  viewsize = 400, 200;
}

ambient_light {
  color = (1, 1, 1)  
}

point_light {
	position = (-5, 5, -5);
	color = (.8, .8, .8);
}

directional_light {
  direction = (-1, -1, 1)  
  color = (0.2, 0.2, 0.2)
}

material {
  name = "wall";
  pattern {
    translate(0, 0.1, 0, scale(0.7,
      checker (0.7, 0.7, 0.7), (0.9, 0.9, 0.9)
    ))
  }
  specular = (0, 0, 0)
}

rotate_x( pi/2, translate( 0, 5, 0,
  plane {
    name = "back";
    material = "wall";
  }
))

translate( 0, -1, 0,
  plane {
    name = "floor";
    material = "wall";
  }
)

translate(2, 0, 0,
  union {
    scale( 0.8, 
      sphere {
        color = "orange"
      }
    ),

    diff {
      sphere {},

      scale( 2, translate(0, 1.2, -0.4,
        sphere {}
      ))
    }
  }
)

translate(-2, 0, 0,
  group {
    intersect {
      scale( 1.3,
        sphere {
          color = "turquoise"
        }
      ),

      cube {
        color = "coral"
      }
    }
  }
)
  `
	scene, err := ParseSbtSceneFromString(data)

	if err == nil {
		scene.World.RenderToPNG(scene.Camera, "test_sbt_csg.png")
	} else {
		t.Errorf("CSG scene failed: %s", err)
	}
}

func TestSbtMeshCow(t *testing.T) {
	TestWithImage(t)

	data := `
FUN-raytracer 0.1

pragma = "gamma=1.0";

camera {
  position = ( 0, 2, -8 );
  target = (0, 0.3, 0);
  updir = ( 0, 1, 0 );
  fovrad = 1;
  viewsize = 400, 200;
}

ambient_light {
  color = (1, 1, 1)  
}

point_light {
	position = (-50, 100, -50);
	color = (.9, .9, .9);
}

rotate_x( pi/2, translate( 0, 5, 0,
  plane { 
    name = "back";
    color = "LightSkyBlue";
  }
))

translate( 0, -1, 0,
  plane {
    name = "floor";
    color = "#567d46";
  }
)

translate( 0, 0.5, -1, scale(2.5,
  group {
    polymesh {
      objfile = "./scenes/cow.obj";
      gennormals = true;
    }
  }
))
  `
	scene, err := ParseSbtSceneFromString(data)

	if err == nil {
		scene.World.RenderToPNG(scene.Camera, "test_sbt_cow.png")
	} else {
		t.Errorf("teapot scene failed: %s", err)
	}
}

func TestSbtDice(t *testing.T) {
	TestWithImage(t)

	data := `
FUN-raytracer 0.1

camera {
  position = ( 0, 2, -9 );
  target = (0, 1, 0);
  updir = ( 0, 1, 0 );
  fovrad = 1.5;
  viewsize = 800, 400;
}

ambient_light {
  color = (1, 1, 1)
}

point_light {
	position = (-5, 5, -3);
	color = (.5, .5, .5);
}

point_light {
	position = (0, 4, -8);
	color = (.4, .4, .4);
}

rotate_x( pi/2, translate( 0, 5, 0,
  plane {
    material = {
      pattern {
        translate(0, 0.1, 0, scale(0.7,
          checker (0.7, 0.7, 0.7), (0.9, 0.9, 0.9)
        ))
      }
      specular = (0, 0, 0)
    }
  }
))

translate( 0, -1, 0,
  plane {
    material = {
      pattern {
        translate(0, 0.1, 0, scale(0.7,
          checker (0.7, 0.7, 0.7), (0.9, 0.9, 0.9)
        ))
      }
      specular = (0.4, 0.4, 0.4)
    }
  }
)

translate(0.5, 2, 1, rotate_x(0.5, rotate_y(5.5, 
  diff {
    group {
      name = "body";

      // Edges
      translate(-1, 0, -1, scale(0.18, 1, 0.18, cyl {} ))
      translate(-1, 0, +1, scale(0.18, 1, 0.18, cyl {} ))
      translate(+1, 0, -1, scale(0.18, 1, 0.18, cyl {} ))
      translate(+1, 0, +1, scale(0.18, 1, 0.18, cyl {} ))

      translate(-1, -1, 0, rotate_x(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(-1, +1, 0, rotate_x(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(+1, -1, 0, rotate_x(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(+1, +1, 0, rotate_x(pi/2, scale(0.18, 1, 0.18, cyl {} )))

      translate(0, -1, -1, rotate_z(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(0, -1, +1, rotate_z(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(0, +1, -1, rotate_z(pi/2, scale(0.18, 1, 0.18, cyl {} )))
      translate(0, +1, +1, rotate_z(pi/2, scale(0.18, 1, 0.18, cyl {} )))

      // Corners
      translate(-1, -1, -1, scale(0.18, sphere {} ))
      translate(-1, +1, -1, scale(0.18, sphere {} ))
      translate(+1, -1, -1, scale(0.18, sphere {} ))
      translate(+1, +1, -1, scale(0.18, sphere {} ))
      translate(-1, -1, +1, scale(0.18, sphere {} ))
      translate(-1, +1, +1, scale(0.18, sphere {} ))
      translate(+1, -1, +1, scale(0.18, sphere {} ))
      translate(+1, +1, +1, scale(0.18, sphere {} ))

      // Faces
      scale(1.18, 1, 1, cube {} )
      scale(1, 1.18, 1, cube {} )
      scale(1, 1, 1.18, cube {} )

      material = {
        diffuse = (1, 0.05, 0) * 0.9;
      }
    },

    group {
      name = "dots";

      translate(0.00, 0.00, -1.18, scale(0.18, sphere {} ))
      translate(1.18, -0.60, -0.60, scale(0.18, sphere {} ))
      translate(1.18, 0.60, 0.60, scale(0.18, sphere {} ))
      translate(0.00, -1.18, 0.00, scale(0.18, sphere {} ))
      translate(-0.60, -1.18, -0.60, scale(0.18, sphere {} ))
      translate(0.60, -1.18, 0.60, scale(0.18, sphere {} ))
      translate(-0.60, 1.18, -0.60, scale(0.18, sphere {} ))
      translate(-0.60, 1.18, 0.60, scale(0.18, sphere {} ))
      translate(0.60, 1.18, -0.60, scale(0.18, sphere {} ))
      translate(0.60, 1.18, 0.60, scale(0.18, sphere {} ))
      translate(-1.18, 0.60, 0.60, scale(0.18, sphere {} ))
      translate(-1.18, -0.60, 0.60, scale(0.18, sphere {} ))
      translate(-1.18, -0.60, -0.60, scale(0.18, sphere {} ))
      translate(-1.18, 0.60, -0.60, scale(0.18, sphere {} ))
      translate(-1.18, 0.00, 0.00, scale(0.18, sphere {} ))
      translate(-0.60, -0.60, 1.18, scale(0.18, sphere {} ))
      translate(-0.60, 0.00, 1.18, scale(0.18, sphere {} ))
      translate(-0.60, 0.60, 1.18, scale(0.18, sphere {} ))
      translate(0.60, -0.60, 1.18, scale(0.18, sphere {} ))
      translate(0.60, 0.00, 1.18, scale(0.18, sphere {} ))
      translate(0.60, 0.60, 1.18, scale(0.18, sphere {} ))  
    }
  }
)))

translate(4, 1, 0, rotate_x(4,
  diff {
    clone "body" {
      material = {
        diffuse = (0, 0.2, 1) * 0.9;
      }
    },

    clone "dots" {
    }
  }
))

translate(-3, 1.5, -1, rotate_x(1, rotate_z(1,
  diff {
    clone "body" {
      material = {
        diffuse = (0.1, 0.8, 0.1) * 0.9;
      }
    },

    clone "dots" {
    }
  }
)))
  `
	scene, err := ParseSbtSceneFromString(data)

	if err == nil {
		scene.World.RenderToPNG(scene.Camera, "test_sbt_dice.png")
	} else {
		t.Errorf("teapot scene failed: %s", err)
	}
}
