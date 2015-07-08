#include <tunables/global>

profile straitjacket/compiler/java {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /usr/lib/jvm/java-7-openjdk-amd64/bin/javac rix,
  /etc/java-7-openjdk/jvm-amd64.cfg r,
  /var/local/straitjacket/tmp/source/?*/?* rw,
}
