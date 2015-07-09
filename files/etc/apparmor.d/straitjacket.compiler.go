#include <tunables/global>

profile straitjacket/compiler/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-compiler>

  /usr/share/go/** rix,
  /usr/lib/go/** rix,
  /usr/src/go/bin/go rix,
  /usr/src/go/src/** r,
  /usr/src/go/pkg/** rix,
  /tmp/** rw,
}
