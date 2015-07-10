#include <tunables/global>

profile straitjacket/compiler/scala {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>
  #include <abstractions/straitjacket-jvm>
  #include <abstractions/straitjacket-compiler>

  /bin/uname rix,
  /usr/lib/jvm/*/bin/java rix,
  /usr/share/scala/lib/** r,
}
