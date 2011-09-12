#include <tunables/global>

profile straitjacket/compiler/scala {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /usr/bin/scalac r,
  /var/local/straitjacket/tmp/source/?*/?* rw,

}
