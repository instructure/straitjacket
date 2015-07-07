#include <tunables/global>

profile straitjacket/compiler/d {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-gcc>

  /var/local/straitjacket/tmp/source/?*/?* r,
  /var/local/straitjacket/tmp/compiler/?* rw,

  /usr/local/bin/gdc rix,
  /usr/bin/gcc* rix,
  /etc/dmd.conf r,
  /tmp/** rw,
  /usr/local/lib64/** r,
  /usr/local/lib/gcc/** rix,
  /usr/local/libexec/gcc/** rix,
  /usr/local/include/d/** r,
}
