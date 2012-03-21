#include <tunables/global>

profile straitjacket/compiled/java {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /var/local/straitjacket/tmp/source/?*/?* r,
}
