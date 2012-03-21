#include <tunables/global>

profile straitjacket/compiler/java {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /etc/java-6-openjdk/jvm.cfg r,
  /usr/lib/jvm/*/bin/javac rix,
  /var/local/straitjacket/tmp/source/?*/?* rw,
}
