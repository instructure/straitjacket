#include <tunables/global>

profile straitjacket/compiler/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-gcc>

  set rlimit nproc <= 32,

  /var/local/straitjacket/tmp/source/?*/?* r,
  /var/local/straitjacket/tmp/source/?*/?*.o rw,
  /var/local/straitjacket/tmp/compiler/?** rw,

  /usr/share/go/** rix,
  /usr/lib/go/** rix,
}
