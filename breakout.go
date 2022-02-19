package main

import (
  "fmt"
  "math/rand"
  "strconv"
)

type PersonCounts struct {
  // Key is "person1-person2" - bool is if they've met yet.
  personToGroupMembers map[string]bool
}

func NewPersonCounts() *PersonCounts {
  m := make(map[string]bool)
  return &PersonCounts{
    personToGroupMembers: m,
  }
}

func (pc *PersonCounts) clone() *PersonCounts {
  newMap := make(map[string]bool, len(pc.personToGroupMembers))
  for k, v := range pc.personToGroupMembers {
    newMap[k] = v
  }
  return &PersonCounts{
    personToGroupMembers: newMap,
  }
}

func (pc *PersonCounts) makeKey(p1, p2 string) string {
  if p1 < p2 {
    return p1 + p2
  }
  return p2 + p1
}

func (pc *PersonCounts) findNewMember(currentGroup []string, remainingPeople []string) (string, bool) {
  bestMember := ""
  // # of new links added
  bestMemberScore := 0
  for _, newP := range remainingPeople {
    score := 0
    for _, groupP := range currentGroup {
      key := pc.makeKey(newP, groupP)
      _, found := pc.personToGroupMembers[key]
      if !found {
        score += 1
      }
    }
    if score > bestMemberScore {
      bestMember = newP
      bestMemberScore = score
    }
  }
  if bestMember != "" {
    for  _, groupP := range currentGroup {
      pc.personToGroupMembers[pc.makeKey(bestMember, groupP)] = true
    }
    // MOVE1: moved up because slice-by-value
    // Could do this instead of order doesn't matter (it probably doesn't but want to maintain for debugging for now
    // remainingPeople[bestMemberIndex] = remainingPeople[len(remainingPeople)-1]
    // return remainingPeople[:len(remainingPeople)-1]
    // remainingPeople = append(remainingPeople[:bestMemberIndex], remainingPeople[bestMemberIndex+1:]...)
  }
  if bestMember == "" {
    // Pick random to avoid person sitting my themselves.
    bestMember = remainingPeople[rand.Intn(len(remainingPeople))]
  }
  return bestMember, bestMemberScore > 0
}

func (pc *PersonCounts) allMatched(people []string) bool {
  for _, p1 := range people {
    for _, p2 := range people {
      if p1 == p2 {
        continue
      }
      _, found := pc.personToGroupMembers[pc.makeKey(p1, p2)]
      if !found {
        return false
      }
    }
  }
  return true
}

type group []string

type breakout []group

func personInBreakout(person string, b breakout) bool {
  for _, g := range b {
    for _, p := range g {
      if person == p {
        return true
      }
    }
  }
  return false
}

func removeFromList(person string, people []string) []string {
  for i, p := range people {
    if p == person {
      return append(people[:i], people[i+1:]...)
    }
  }
  return people
}

func generatePeople(length int) []string {
  var people []string
  for i := 0; i < length; i++ {
    people = append(people, "person_" + strconv.Itoa(i))
  }
  return people
}

func generateGroup(people []string, groupSize int, b breakout, pc *PersonCounts) (group, bool) {
  peopleRemaining := make([]string, len(people))
  copy(peopleRemaining, people)

  var g group
  r := rand.Intn(len(peopleRemaining))
  g = append(g, peopleRemaining[r])
  peopleRemaining = append(peopleRemaining[:r], peopleRemaining[r+1:]...)

  anyAddLinks := false
  for i := 1; i < groupSize; i++ {
    newP, addsLinks := pc.findNewMember(g, peopleRemaining)
    peopleRemaining = removeFromList(newP, peopleRemaining)
    if addsLinks {
      anyAddLinks = true
    }
    g = append(g, newP)
  }
  return g, anyAddLinks
}

func generateBreakout(people []string, groupSize int, breakouts []breakout, pc *PersonCounts) breakout {
  peopleRemaining := make([]string, len(people))
  copy(peopleRemaining, people)
  pcClone := pc.clone()

  numGroups := len(people) / groupSize
  var b breakout
  anyAddLinks := false
  for i := 0; i < numGroups; i++ {
    g, addLinks := generateGroup(peopleRemaining, groupSize, b, pc)
    for _, p := range g {
      peopleRemaining = removeFromList(p, peopleRemaining)
    }
    b = append(b, g)
    if addLinks {
      anyAddLinks = true
    }
  }
  if anyAddLinks {
    return b
  }
  // This is bad but it works.
  // Recurse until allMatched is true.
  if !pc.allMatched(people) {
    return generateBreakout(people, groupSize, breakouts, pcClone)
  }
  return nil
}

func main() {
  // TODO: rounding for when these are not divisible
  numPeople := 12
  groupSize := 3
  people := generatePeople(numPeople)

  var breakouts []breakout
  pc := NewPersonCounts()
  for {
    newBreakout := generateBreakout(people, groupSize, breakouts, pc)
    if newBreakout == nil {
      break
    }
    breakouts = append(breakouts, newBreakout)
  }

  for i, breakout := range breakouts {
    fmt.Printf("breakout %d:\n", i)
    for j, g := range breakout {
      fmt.Printf("\tgroup %d: %v\n", j, g)
    }
  }
}
