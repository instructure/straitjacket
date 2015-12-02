#include <tunables/global>

profile straitjacket/interpreter/ruby {
  #include <abstractions/base>
  #include <abstractions/ruby>
  #include <abstractions/straitjacket-base>

# what the heck ruby
  set rlimit nproc <= 300,

  /usr/local/lib/ruby/*/** r,
  /usr/local/lib/ruby/*/**.so rm,
}
