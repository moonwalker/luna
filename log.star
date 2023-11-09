def _info(s):
  print('>>>', s)

log = module(
  "log",
  info = _info,
)
