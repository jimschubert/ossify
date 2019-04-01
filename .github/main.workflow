workflow "Master Build" {
  on = "push"
  resolves = ["ossify github action"]
}

action "ossify github action" {
  uses = "./ci"
}
