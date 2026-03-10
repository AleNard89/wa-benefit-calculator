export type Conversation = {
  id: number
  userId: number
  companyId: number
  title: string
  createdAt: string
  updatedAt: string
}

export type Message = {
  id: number
  conversationId: number
  role: 'user' | 'assistant'
  content: string
  createdAt: string
}
