#include <tunables/global>

profile straitjacket/compiled/java {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /etc/java-7-openjdk/jvm-amd64.cfg r,
  /var/local/straitjacket/tmp/source/?*/?* rw,
}
