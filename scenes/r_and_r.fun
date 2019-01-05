FUN-raytracer 1.0

// Initial scene from book's chapter 11
// Author: Jamis Buck

pragma = "gamma=1.0";

camera {
  position = ( -2.6, 1.5, -3.9 );
  target = (-0.6, 1, -0.8);
  updir = ( 0, 1, 0 );
  fovrad = 1.152;
}

ambient_light {
  color = (1, 1, 1)
}

point_light {
  position = (-4.9, 4.9, -1);
  color = (1, 1, 1);
}

material {
  name = "wall";
  pattern {
    scale(0.25, rotate(0, 1, 0, 1.5708,
      stripe (0.45, 0.45, 0.45), (0.55, 0.55, 0.55)
    ))
  }
  ambient = (0, 0, 0)
  diffuse = (0.4, 0.4, 0.4)
  specular = (0, 0, 0)
  reflective = (0.3, 0.3, 0.3)
}

rotate( 0, 1, 0, 0.31415,
  plane {
    name = "floor";
    material = {
      pattern {
        translate(0, 0.1, 0,
          checker (0.65, 0.65, 0.65), (0.35, 0.35, 0.35)
        )
      }
      specular = (0, 0, 0)
      reflective = (0.4, 0.4, 0.4)
    }
  }
)

translate(0, 5, 0,
  plane {
    name = "ceiling";
    material = {
      diffuse = (0.8, 0.8, 0.8)
      specular = (0, 0, 0)
      ambient = (0.3, 0.3, 0.3)
    }
  }
)

translate(-5, 0, 0, rotate(0, 0, 1, 1.5708, rotate(0, 1, 0, 1.5708,
  plane {
    name = "west wall";
    material = "wall";
  }
)))

translate(5, 0, 0, rotate(0, 0, 1, 1.5708, rotate(0, 1, 0, 1.5708,
  plane {
    name = "east wall";
    material = "wall";
  }
)))

translate(0, 0, 5, rotate(1, 0, 0, 1.5708,
  plane {
    name = "north wall";
    material = "wall";
  }
))

translate(0, 0, -5, rotate(1, 0, 0, 1.5708,
  plane {
    name = "south wall";
    material = "wall";
  }
))

translate(4.6, 0.4, 1, scale( 0.4,
  sphere {
    material = {
      diffuse = (0.8, 0.5, 0.3) * 0.9;
      shininess = 50;
    }
  }
))

translate(4.7, 0.3, 0.4, scale( 0.3,
  sphere {
    material = {
      diffuse = (0.9, 0.4, 0.5) * 0.9;
      shininess = 50;
    }
  }
))

translate(-1, 0.5, 4.5, scale( 0.5,
  sphere {
    material = {
      diffuse = (0.4, 0.9, 0.6) * 0.9;
      shininess = 50;
    }
  }
))

translate(-1.7, 0.3, 4.7, scale( 0.3,
  sphere {
    material = {
      diffuse = (0.4, 0.6, 0.9) * 0.9;
      shininess = 50;
    }
  }
))

translate(-0.6, 1, 0.6, scale( 1,
  sphere {
    material = {
      diffuse = (1, 0.3, 0.2) * 0.9;
      specular = (0.4, 0.4, 0.4);
      shininess = 5;
    }
  }
))

translate(0.6, 0.7, -0.6, scale( 0.7,
  sphere {
    material = {
      diffuse = (0, 0, 0.2) * 0.4;
      ambient = (0, 0, 0);
      specular = (0.9, 0.9, 0.9);
      shininess = 300;
      reflective = (0.9, 0.9, 0.9);
      transmissive = (0.9, 0.9, 0.9);
      index = 1.5;
    }
  }
))

translate(-0.7, 0.5, -0.8, scale( 0.5,
  sphere {
    material = {
      diffuse = (0, 0.2, 0) * 0.4;
      ambient = (0, 0, 0);
      specular = (0.9, 0.9, 0.9);
      shininess = 300;
      reflective = (0.9, 0.9, 0.9);
      transmissive = (0.9, 0.9, 0.9);
      index = 1.5;
    }
  }
))
