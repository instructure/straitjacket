#include <tunables/global>

profile straitjacket/compiler/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-compiler>

  /usr/local/go/** rix,
  /tmp/** rw,
}
