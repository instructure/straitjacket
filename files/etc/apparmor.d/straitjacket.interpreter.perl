#include <tunables/global>

profile straitjacket/interpreter/perl {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>

  /usr/share/perl/**/** rix,
  /usr/lib/*/perl/** rix,
}
