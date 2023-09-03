FUN-raytracer 0.1

// A glass sphere with a spherical hole inside, demostrating refraction and the Fresnel effect
// (it is recommended to render it with at least 16 samples per pixel)

pragma = "gamma=1.0";

translate(0, -10.1, 0,
  plane {
    material = {
      pattern {
        translate(0, 0.1, 0,
          checker "black", "white"
        )
      }
    }
  }
)

sphere {
  material = {
    diffuse = (1,1,1) * 0.1;
    shininess = 300;
    reflective = (1, 1, 1);
    transmissive = (1, 1, 1);
    index = 1.52;
  }
}

scale(0.5,
  sphere {
    material = {
      diffuse = (1,1,1) * 0.1;
      shininess = 300;
      reflective = (1, 1, 1);
      transmissive = (1, 1, 1);
      index = 1;
    }
  }
)

camera {
  position = ( 0, 2.5, 0 );
  target = (0, 0, 0);
  updir = ( 1, 0, 0 );
  fovrad = pi/3;
  viewsize = 480, 480;
}

ambient_light {
  color = (1, 1, 1)
}

point_light {
	position = (20, 10, 0);
	color = (.7, .7, .7);
}

