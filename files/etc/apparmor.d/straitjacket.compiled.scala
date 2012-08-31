#include <tunables/global>

profile straitjacket/compiled/scala {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /bin/uname rix,
  /usr/bin/scala rix,
  /usr/lib/jvm/*/bin/java rix,
  /var/local/straitjacket/tmp/source/?*/?* r,

}
