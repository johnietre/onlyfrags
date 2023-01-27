package main

// Item represents an item for sale
type Item struct {
  // Path to parent 
  ParentImg string
  // Paths to images associated with the item
  Imgs []string
}
