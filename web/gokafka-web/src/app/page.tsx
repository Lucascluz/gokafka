import Link from "next/link"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

// Mock data for products
const featuredProducts = [
  {
    id: 1,
    title: "iPhone 14 Pro",
    price: 899,
    image: "/placeholder-phone.jpg",
    location: "New York, NY",
    category: "Electronics",
  },
  {
    id: 2,
    title: "Nike Air Max 90",
    price: 120,
    image: "/placeholder-shoes.jpg",
    location: "Los Angeles, CA",
    category: "Fashion",
  },
  {
    id: 3,
    title: "MacBook Pro 2023",
    price: 1299,
    image: "/placeholder-laptop.jpg",
    location: "San Francisco, CA",
    category: "Electronics",
  },
  {
    id: 4,
    title: "Vintage Leather Jacket",
    price: 85,
    image: "/placeholder-jacket.jpg",
    location: "Chicago, IL",
    category: "Fashion",
  },
]

const categories = [
  { name: "Electronics", count: "2,847 ads" },
  { name: "Fashion", count: "1,923 ads" },
  { name: "Home & Garden", count: "1,456 ads" },
  { name: "Cars", count: "892 ads" },
  { name: "Sports", count: "734 ads" },
  { name: "Books", count: "523 ads" },
]

export default function HomePage() {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <section className="text-center py-12 mb-12">
        <h1 className="text-4xl font-bold tracking-tight sm:text-6xl mb-6">
          Find everything you need
        </h1>
        <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
          Buy and sell locally. Discover amazing deals on products from people near you.
        </p>
        <div className="flex gap-4 justify-center">
          <Button size="lg" asChild>
            <Link href="/products">Browse Products</Link>
          </Button>
          <Button variant="outline" size="lg" asChild>
            <Link href="/products/new">Start Selling</Link>
          </Button>
        </div>
      </section>

      {/* Categories */}
      <section className="mb-12">
        <h2 className="text-2xl font-bold mb-6">Browse by Category</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
          {categories.map((category) => (
            <Card key={category.name} className="hover:shadow-md transition-shadow cursor-pointer">
              <CardContent className="p-4 text-center">
                <h3 className="font-medium mb-1">{category.name}</h3>
                <p className="text-sm text-muted-foreground">{category.count}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      {/* Featured Products */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">Featured Products</h2>
          <Button variant="outline" asChild>
            <Link href="/products">View All</Link>
          </Button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {featuredProducts.map((product) => (
            <Card key={product.id} className="hover:shadow-lg transition-shadow cursor-pointer">
              <div className="aspect-square bg-muted rounded-t-lg"></div>
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <Badge variant="secondary">{product.category}</Badge>
                  <span className="text-lg font-bold">${product.price}</span>
                </div>
                <CardTitle className="text-lg">{product.title}</CardTitle>
                <CardDescription>{product.location}</CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </section>

      {/* CTA Section */}
      <section className="text-center py-12 mt-16 bg-muted rounded-lg">
        <h2 className="text-3xl font-bold mb-4">Ready to start selling?</h2>
        <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
          Join thousands of sellers already making money on GoMarket. It&apos;s free and easy to get started.
        </p>
        <Button size="lg" asChild>
          <Link href="/register">Get Started</Link>
        </Button>
      </section>
    </div>
  )
}
