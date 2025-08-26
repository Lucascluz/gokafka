"use client"

import { useState } from "react"
import Link from "next/link"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Plus, MoreVertical, Edit, Trash2, Eye, PauseCircle, PlayCircle } from "lucide-react"

// Mock data for user's products
const userProducts = [
  {
    id: 1,
    title: "iPhone 14 Pro - Excellent Condition",
    price: 899,
    image: "/placeholder-phone.jpg",
    category: "Electronics",
    condition: "Like New",
    status: "active",
    views: 156,
    favorites: 23,
    messages: 8,
    postedDate: "2024-01-15",
    description: "Selling my iPhone 14 Pro in excellent condition. Barely used, always kept in a case with screen protector.",
  },
  {
    id: 2,
    title: "Gaming Setup - Complete Package",
    price: 1200,
    image: "/placeholder-setup.jpg",
    category: "Electronics",
    condition: "Very Good",
    status: "active",
    views: 89,
    favorites: 12,
    messages: 5,
    postedDate: "2024-01-10",
    description: "Complete gaming setup including monitor, keyboard, mouse, and headset.",
  },
  {
    id: 3,
    title: "Vintage Leather Jacket",
    price: 85,
    image: "/placeholder-jacket.jpg",
    category: "Fashion",
    condition: "Good",
    status: "paused",
    views: 34,
    favorites: 7,
    messages: 2,
    postedDate: "2024-01-05",
    description: "Authentic vintage leather jacket from the 80s. Size medium.",
  },
  {
    id: 4,
    title: "Road Bike - Trek",
    price: 450,
    image: "/placeholder-bike.jpg",
    category: "Sports",
    condition: "Good",
    status: "sold",
    views: 78,
    favorites: 15,
    messages: 12,
    postedDate: "2023-12-28",
    soldDate: "2024-01-08",
    description: "Well-maintained Trek road bike. Perfect for commuting or weekend rides.",
  },
]

const getStatusColor = (status: string) => {
  switch (status) {
    case "active":
      return "bg-green-100 text-green-800"
    case "paused":
      return "bg-yellow-100 text-yellow-800"
    case "sold":
      return "bg-blue-100 text-blue-800"
    default:
      return "bg-gray-100 text-gray-800"
  }
}

export default function MyProductsPage() {
  const [products, setProducts] = useState(userProducts)
  const [activeTab, setActiveTab] = useState("all")

  const handleStatusChange = (productId: number, newStatus: string) => {
    setProducts(products.map(product => 
      product.id === productId ? { ...product, status: newStatus } : product
    ))
  }

  const handleDelete = (productId: number) => {
    if (confirm("Are you sure you want to delete this listing?")) {
      setProducts(products.filter(product => product.id !== productId))
    }
  }

  const filteredProducts = products.filter(product => {
    if (activeTab === "all") return true
    return product.status === activeTab
  })

  const stats = {
    total: products.length,
    active: products.filter(p => p.status === "active").length,
    paused: products.filter(p => p.status === "paused").length,
    sold: products.filter(p => p.status === "sold").length,
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">My Products</h1>
            <p className="text-muted-foreground">Manage your listings and track performance</p>
          </div>
          <Button asChild>
            <Link href="/products/new">
              <Plus className="h-4 w-4 mr-2" />
              Add New Product
            </Link>
          </Button>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4 text-center">
              <div className="text-2xl font-bold">{stats.total}</div>
              <div className="text-sm text-muted-foreground">Total Listings</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="text-2xl font-bold text-green-600">{stats.active}</div>
              <div className="text-sm text-muted-foreground">Active</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="text-2xl font-bold text-yellow-600">{stats.paused}</div>
              <div className="text-sm text-muted-foreground">Paused</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="text-2xl font-bold text-blue-600">{stats.sold}</div>
              <div className="text-sm text-muted-foreground">Sold</div>
            </CardContent>
          </Card>
        </div>

        {/* Products List */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="all">All ({stats.total})</TabsTrigger>
            <TabsTrigger value="active">Active ({stats.active})</TabsTrigger>
            <TabsTrigger value="paused">Paused ({stats.paused})</TabsTrigger>
            <TabsTrigger value="sold">Sold ({stats.sold})</TabsTrigger>
          </TabsList>

          <TabsContent value={activeTab} className="space-y-4">
            {filteredProducts.length === 0 ? (
              <Card>
                <CardContent className="p-12 text-center">
                  <h3 className="text-lg font-medium mb-2">No products found</h3>
                  <p className="text-muted-foreground mb-4">
                    {activeTab === "all" 
                      ? "You haven't listed any products yet." 
                      : `You don't have any ${activeTab} products.`}
                  </p>
                  <Button asChild>
                    <Link href="/products/new">List Your First Product</Link>
                  </Button>
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-4">
                {filteredProducts.map((product) => (
                  <Card key={product.id}>
                    <CardContent className="p-6">
                      <div className="flex gap-4">
                        <div className="w-20 h-20 bg-muted rounded-lg flex-shrink-0"></div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-start justify-between">
                            <div className="flex-1 min-w-0">
                              <h3 className="font-semibold text-lg truncate">{product.title}</h3>
                              <p className="text-sm text-muted-foreground mb-2 line-clamp-2">
                                {product.description}
                              </p>
                              <div className="flex items-center gap-4 text-sm text-muted-foreground">
                                <span>Posted: {new Date(product.postedDate).toLocaleDateString()}</span>
                                <span>•</span>
                                <span>{product.category}</span>
                                <span>•</span>
                                <span>{product.condition}</span>
                              </div>
                            </div>
                            <div className="flex items-center gap-4 ml-4">
                              <div className="text-right">
                                <div className="text-xl font-bold text-green-600">${product.price}</div>
                                <Badge className={getStatusColor(product.status)}>
                                  {product.status}
                                </Badge>
                              </div>
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="icon">
                                    <MoreVertical className="h-4 w-4" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem asChild>
                                    <Link href={`/products/${product.id}`}>
                                      <Eye className="h-4 w-4 mr-2" />
                                      View
                                    </Link>
                                  </DropdownMenuItem>
                                  <DropdownMenuItem asChild>
                                    <Link href={`/products/${product.id}/edit`}>
                                      <Edit className="h-4 w-4 mr-2" />
                                      Edit
                                    </Link>
                                  </DropdownMenuItem>
                                  {product.status === "active" ? (
                                    <DropdownMenuItem 
                                      onClick={() => handleStatusChange(product.id, "paused")}
                                    >
                                      <PauseCircle className="h-4 w-4 mr-2" />
                                      Pause
                                    </DropdownMenuItem>
                                  ) : product.status === "paused" ? (
                                    <DropdownMenuItem 
                                      onClick={() => handleStatusChange(product.id, "active")}
                                    >
                                      <PlayCircle className="h-4 w-4 mr-2" />
                                      Activate
                                    </DropdownMenuItem>
                                  ) : null}
                                  <DropdownMenuItem 
                                    onClick={() => handleDelete(product.id)}
                                    className="text-red-600"
                                  >
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </div>
                          
                          {/* Stats */}
                          <div className="flex items-center gap-6 mt-4 text-sm text-muted-foreground">
                            <div className="flex items-center gap-1">
                              <Eye className="h-4 w-4" />
                              {product.views} views
                            </div>
                            <div className="flex items-center gap-1">
                              <span>❤️</span>
                              {product.favorites} favorites
                            </div>
                            <div className="flex items-center gap-1">
                              <span>💬</span>
                              {product.messages} messages
                            </div>
                            {product.status === "sold" && product.soldDate && (
                              <div>
                                Sold on {new Date(product.soldDate).toLocaleDateString()}
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
