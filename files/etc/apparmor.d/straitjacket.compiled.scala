#include <tunables/global>

profile straitjacket/compiled/scala {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /usr/bin/scala r,
  /var/local/straitjacket/tmp/source/?*/?* r,

}
