"use client"

import { useState } from "react"
import Link from "next/link"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Search, MapPin, Heart } from "lucide-react"

// Mock data for products
const products = [
  {
    id: 1,
    title: "iPhone 14 Pro - Excellent Condition",
    price: 899,
    image: "/placeholder-phone.jpg",
    location: "New York, NY",
    category: "Electronics",
    condition: "Like New",
    seller: "John D.",
    postedDate: "2 days ago",
    isFeatured: true,
  },
  {
    id: 2,
    title: "Nike Air Max 90 - Size 10",
    price: 120,
    image: "/placeholder-shoes.jpg",
    location: "Los Angeles, CA",
    category: "Fashion",
    condition: "Good",
    seller: "Sarah M.",
    postedDate: "1 week ago",
    isFeatured: false,
  },
  {
    id: 3,
    title: "MacBook Pro 2023 - 16GB RAM",
    price: 1299,
    image: "/placeholder-laptop.jpg",
    location: "San Francisco, CA",
    category: "Electronics",
    condition: "Excellent",
    seller: "Mike R.",
    postedDate: "3 days ago",
    isFeatured: true,
  },
  {
    id: 4,
    title: "Vintage Leather Jacket - Medium",
    price: 85,
    image: "/placeholder-jacket.jpg",
    location: "Chicago, IL",
    category: "Fashion",
    condition: "Good",
    seller: "Lisa K.",
    postedDate: "5 days ago",
    isFeatured: false,
  },
  {
    id: 5,
    title: "Professional Camera Kit",
    price: 450,
    image: "/placeholder-camera.jpg",
    location: "Austin, TX",
    category: "Electronics",
    condition: "Very Good",
    seller: "Alex P.",
    postedDate: "1 day ago",
    isFeatured: false,
  },
  {
    id: 6,
    title: "Gaming Chair - Almost New",
    price: 200,
    image: "/placeholder-chair.jpg",
    location: "Seattle, WA",
    category: "Home & Garden",
    condition: "Like New",
    seller: "Chris B.",
    postedDate: "4 days ago",
    isFeatured: false,
  },
]

const categories = ["All", "Electronics", "Fashion", "Home & Garden", "Cars", "Sports", "Books"]
const sortOptions = ["Newest", "Price: Low to High", "Price: High to Low", "Distance"]

export default function ProductsPage() {
  const [searchTerm, setSearchTerm] = useState("")
  const [selectedCategory, setSelectedCategory] = useState("All")
  const [sortBy, setSortBy] = useState("Newest")
  const [filteredProducts, setFilteredProducts] = useState(products)

  const handleSearch = (term: string) => {
    setSearchTerm(term)
    filterProducts(term, selectedCategory)
  }

  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category)
    filterProducts(searchTerm, category)
  }

  const filterProducts = (term: string, category: string) => {
    let filtered = products

    if (term) {
      filtered = filtered.filter(product =>
        product.title.toLowerCase().includes(term.toLowerCase()) ||
        product.category.toLowerCase().includes(term.toLowerCase())
      )
    }

    if (category !== "All") {
      filtered = filtered.filter(product => product.category === category)
    }

    setFilteredProducts(filtered)
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Browse Products</h1>
            <p className="text-muted-foreground">
              Discover amazing deals from people near you
            </p>
          </div>
          <Button asChild>
            <Link href="/products/new">Sell Item</Link>
          </Button>
        </div>

        {/* Search and Filters */}
        <div className="flex flex-col md:flex-row gap-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search products..."
              className="pl-10"
              value={searchTerm}
              onChange={(e) => handleSearch(e.target.value)}
            />
          </div>
          <div className="flex gap-2">
            <Select value={selectedCategory} onValueChange={handleCategoryChange}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                {categories.map((category) => (
                  <SelectItem key={category} value={category}>
                    {category}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={sortBy} onValueChange={setSortBy}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                {sortOptions.map((option) => (
                  <SelectItem key={option} value={option}>
                    {option}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Results Count */}
        <div className="text-sm text-muted-foreground">
          Showing {filteredProducts.length} of {products.length} products
        </div>

        {/* Products Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {filteredProducts.map((product) => (
            <Card key={product.id} className="hover:shadow-lg transition-shadow cursor-pointer group">
              <Link href={`/products/${product.id}`}>
                <div className="relative">
                  <div className="aspect-square bg-muted rounded-t-lg"></div>
                  {product.isFeatured && (
                    <Badge className="absolute top-2 left-2" variant="secondary">
                      Featured
                    </Badge>
                  )}
                  <Button
                    variant="ghost"
                    size="icon"
                    className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity"
                    onClick={(e) => {
                      e.preventDefault()
                      // Handle favorite logic
                    }}
                  >
                    <Heart className="h-4 w-4" />
                  </Button>
                </div>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <Badge variant="outline">{product.category}</Badge>
                    <span className="text-lg font-bold text-green-600">${product.price}</span>
                  </div>
                  <CardTitle className="text-lg line-clamp-2">{product.title}</CardTitle>
                  <CardDescription className="space-y-1">
                    <div className="flex items-center text-sm">
                      <MapPin className="h-3 w-3 mr-1" />
                      {product.location}
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Condition: {product.condition}</span>
                      <span className="text-muted-foreground">{product.postedDate}</span>
                    </div>
                  </CardDescription>
                </CardHeader>
                <CardContent className="pt-0">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">
                      by {product.seller}
                    </span>
                    <Button size="sm" variant="outline">
                      View Details
                    </Button>
                  </div>
                </CardContent>
              </Link>
            </Card>
          ))}
        </div>

        {/* Load More */}
        {filteredProducts.length > 0 && (
          <div className="text-center">
            <Button variant="outline" size="lg">
              Load More Products
            </Button>
          </div>
        )}

        {/* No Results */}
        {filteredProducts.length === 0 && (
          <div className="text-center py-12">
            <h3 className="text-lg font-medium mb-2">No products found</h3>
            <p className="text-muted-foreground mb-4">
              Try adjusting your search terms or filters
            </p>
            <Button variant="outline" onClick={() => {
              setSearchTerm("")
              setSelectedCategory("All")
              setFilteredProducts(products)
            }}>
              Clear Filters
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
