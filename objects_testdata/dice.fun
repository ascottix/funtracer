FUN-raytracer 0.1

// Three rolling dice demonstrating Constructive Solid Geometry and cloning

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

      // Set color for all group members
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
