#include <tunables/global>

profile straitjacket/compiler/c {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-compiler>
  #include <abstractions/straitjacket-gcc>

  /tmp/** rw,
}
