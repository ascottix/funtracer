FUN-raytracer 0.1

// A simple wrapper for teapot.obj demonstrating triangular meshes

camera {
  position = ( 0, 2, -8 );
  target = (0, 0, 0);
  updir = ( 0, 1, 0 );
  fovrad = 1;
  viewsize = 400, 200;
}

ambient_light {
  color = (0.5, 0.5, 0.5)
}

point_light {
	position = (-5, 10, -5);
	color = (.9, .9, .9);
}

material {
  name = "wall";
  pattern {
    translate(0, 0.1, 0, scale(0.7,
      checker (0.5, 0.5, 0.5), (0.7, 0.7, 0.7)
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

rotate_x( -pi/2, scale(3,
  group {
    polymesh {
      name = "teapot";
      objfile = "teapot.obj";
      material = {
        diffuse = (0.87, 0.87, 0.9) * 1;
        reflective = (0.3, 0.3, 0.3);
        shininess = 250;
      }
    }
  }
))
