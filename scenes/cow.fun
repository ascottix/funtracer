FUN-raytracer 0.1

// A simple wrapper for cow.obj demonstrating automatic generation of normals
// (try setting gennormals to false)

camera {
  position = ( 0, 2, -8 );
  target = (0, 0.3, 0);
  updir = ( 0, 1, 0 );
  fovrad = 1;
  viewsize = 800, 400;
}

ambient_light {
  color = (0.5,0.5,0.5)
}

point_light {
	position = (-5, 10, -5);
	color = (1, 1, 1);
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
      gennormals = true;  // Generate vertex normals, this option must be declared before the OBJ file
      objfile = "cow.obj";
    }
  }
))
