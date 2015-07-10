#include <tunables/global>

profile straitjacket/compiler/cpp {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-compiler>
  #include <abstractions/straitjacket-gcc>

  /usr/bin/g++* rix,
  /tmp/** rw,
}
