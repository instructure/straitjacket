#include <tunables/global>

profile straitjacket/compiler/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-gcc>

  /var/local/straitjacket/tmp/source/?*/?* r,
  /var/local/straitjacket/tmp/compiler/?* rw,

  /usr/bin/golang-g rix,
  /usr/bin/golang-l rix,
  /usr/bin/6g rix,
  /usr/bin/6l rix,
}
