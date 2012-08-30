#include <tunables/global>

profile straitjacket/compiler/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-gcc>

  /var/local/straitjacket/tmp/source/?*/?* r,
  /var/local/straitjacket/tmp/source/?*/?*.o rw,
  /var/local/straitjacket/tmp/compiler/?** rw,

  /usr/share/go/** rix,
  /usr/lib/go/** rix,
}
