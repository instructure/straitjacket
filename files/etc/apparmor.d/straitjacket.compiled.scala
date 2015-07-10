#include <tunables/global>

profile straitjacket/compiled/scala {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>

  /bin/uname rix,
  /usr/lib/jvm/*/bin/java rix,
  /usr/share/scala/lib/** r,
}
