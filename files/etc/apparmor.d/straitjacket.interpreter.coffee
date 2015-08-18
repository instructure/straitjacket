#include <tunables/global>

profile straitjacket/interpreter/coffee {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>

# what the heck, node loves forking
  set rlimit nproc <= 250,
# node links in a lot
  set rlimit as <= 1G,

  /usr/bin/env rix,
  /usr/local/bin/node rix,
  /usr/local/lib/node_modules/coffee-script/**/* rix,
  /var/local/straitjacket/tmp/source/?*/?* r,
}
