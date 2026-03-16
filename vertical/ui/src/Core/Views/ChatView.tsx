import { Box, Button, Flex, Input, Text } from '@chakra-ui/react'
import { useCallback, useEffect, useRef, useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { LuMessageCircle, LuPlus, LuSend, LuTrash2 } from 'react-icons/lu'

import Config from '@/Config'
import { useCurrentUser } from '@/Auth/Hooks'
import {
  useConversationsQuery,
  useConversationMessagesQuery,
  useCreateConversationMutation,
  useDeleteConversationMutation,
} from '@/Chat/Services/Api'
import { api } from '@/Core/Services/Api'
import type { RootState } from '@/Core/Redux/Store'
import { useDispatch, useSelector } from 'react-redux'
import type { Message } from '@/Chat/Types'

function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' })
}

function MessageBubble({ message }: { message: Message }) {
  const isUser = message.role === 'user'
  return (
    <Flex justify={isUser ? 'flex-end' : 'flex-start'} mb={3}>
      <Box
        maxW="75%"
        px={4}
        py={2.5}
        borderRadius={isUser ? '16px 16px 4px 16px' : '16px 16px 16px 4px'}
        bg={isUser ? '#007aff' : '#f5f5f7'}
        color={isUser ? 'white' : '#1d1d1f'}
      >
        {isUser ? (
          <Text fontSize="14px" lineHeight="1.5" whiteSpace="pre-wrap">{message.content}</Text>
        ) : (
          <Box className="chat-markdown" fontSize="14px" lineHeight="1.5">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{message.content}</ReactMarkdown>
          </Box>
        )}
        <Text fontSize="10px" color={isUser ? 'whiteAlpha.700' : '#86868b'} mt={1} textAlign="right">
          {formatTime(message.createdAt)}
        </Text>
      </Box>
    </Flex>
  )
}

function TypingIndicator() {
  return (
    <Flex justify="flex-start" mb={3}>
      <Box px={4} py={3} borderRadius="16px 16px 16px 4px" bg="#f5f5f7">
        <Flex gap="4px" align="center" h="20px">
          {[0, 1, 2].map((i) => (
            <Box
              key={i}
              w="8px"
              h="8px"
              borderRadius="50%"
              bg="#86868b"
              animation={`typingBounce 1.4s ease-in-out ${i * 0.2}s infinite`}
            />
          ))}
        </Flex>
      </Box>
    </Flex>
  )
}

function StreamingBubble({ content }: { content: string }) {
  if (!content) return null
  return (
    <Flex justify="flex-start" mb={3}>
      <Box maxW="75%" px={4} py={2.5} borderRadius="16px 16px 16px 4px" bg="#f5f5f7" color="#1d1d1f">
        <Box className="chat-markdown" fontSize="14px" lineHeight="1.5">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
          <Box as="span" display="inline-block" w="6px" h="14px" bg="#007aff" borderRadius="1px" ml={0.5} animation="blink 1s infinite" />
        </Box>
      </Box>
    </Flex>
  )
}

export default function ChatView() {
  const user = useCurrentUser()
  const dispatch = useDispatch()
  const token = useSelector((s: RootState) => s.auth.token)
  const companyId = useSelector((s: RootState) => s.orgs.companyId)
  const { data: conversations } = useConversationsQuery()
  const [createConversation] = useCreateConversationMutation()
  const [deleteConversation] = useDeleteConversationMutation()

  const [activeConvId, setActiveConvId] = useState<number | null>(null)
  const [inputValue, setInputValue] = useState('')
  const [streaming, setStreaming] = useState(false)
  const [streamContent, setStreamContent] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const { data: messages } = useConversationMessagesQuery(activeConvId!, { skip: !activeConvId })

  useEffect(() => {
    if (conversations && conversations.length > 0 && !activeConvId) {
      setActiveConvId(conversations[0].id)
    }
  }, [conversations, activeConvId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, streamContent])

  const handleSend = useCallback(async () => {
    if (!inputValue.trim() || streaming) return
    const content = inputValue.trim()
    setInputValue('')

    let convId = activeConvId
    if (!convId) {
      try {
        const conv = await createConversation({ title: 'Nuova conversazione' }).unwrap()
        convId = conv.id
        setActiveConvId(conv.id)
      } catch {
        return
      }
    }

    // Optimistically add user message to cache
    dispatch(
      api.util.updateQueryData('conversationMessages', convId, (draft: Message[]) => {
        draft.push({
          id: Date.now(),
          conversationId: convId!,
          role: 'user',
          content,
          createdAt: new Date().toISOString(),
        })
      }),
    )

    setStreaming(true)
    setStreamContent('')

    try {
      const response = await fetch(`${Config.api.basePath}/chat/conversations/${convId}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
          'X-Company-Id': companyId,
        },
        body: JSON.stringify({ content }),
      })

      if (!response.ok || !response.body) {
        setStreaming(false)
        return
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let fullContent = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const text = decoder.decode(value, { stream: true })
        const lines = text.split('\n')

        for (const line of lines) {
          if (!line.startsWith('data: ')) continue
          const data = line.slice(6)
          if (data === '[DONE]') continue

          try {
            const parsed = JSON.parse(data)
            const delta = parsed.choices?.[0]?.delta?.content
            if (delta) {
              fullContent += delta
              setStreamContent(fullContent)
            }
          } catch {
            // skip malformed chunks
          }
        }
      }

      // Add assistant message to cache
      if (fullContent) {
        dispatch(
          api.util.updateQueryData('conversationMessages', convId, (draft: Message[]) => {
            draft.push({
              id: Date.now() + 1,
              conversationId: convId!,
              role: 'assistant',
              content: fullContent,
              createdAt: new Date().toISOString(),
            })
          }),
        )
      }
    } catch (err) {
      console.error('Stream error', err)
    } finally {
      setStreaming(false)
      setStreamContent('')
      // Refresh conversations list (title may have changed)
      dispatch(api.util.invalidateTags(['Conversations']))
    }
  }, [inputValue, activeConvId, streaming, token, companyId, createConversation, dispatch])

  const handleNewChat = async () => {
    try {
      const conv = await createConversation({ title: 'Nuova conversazione' }).unwrap()
      setActiveConvId(conv.id)
    } catch { /* */ }
  }

  const handleDelete = async (id: number, e: React.MouseEvent) => {
    e.stopPropagation()
    await deleteConversation(id).unwrap()
    if (activeConvId === id) setActiveConvId(null)
  }

  return (
    <Flex h="calc(100vh - 48px)" gap={0}>
      {/* Sidebar - Conversations */}
      <Flex
        direction="column"
        w="260px"
        minW="260px"
        bg="white"
        borderRadius="16px"
        boxShadow="0 2px 20px rgba(0,0,0,0.06)"
        mr={4}
        overflow="hidden"
      >
        <Flex align="center" justify="space-between" px={4} pt={4} pb={3}>
          <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Chat</Text>
          <Box
            as="button"
            p={1.5}
            borderRadius="8px"
            color="#007aff"
            cursor="pointer"
            _hover={{ bg: '#007aff10' }}
            onClick={handleNewChat}
          >
            <LuPlus size={18} />
          </Box>
        </Flex>

        <Flex direction="column" flex={1} overflow="auto" px={2} pb={2} gap={0.5}>
          {conversations && conversations.length > 0 ? (
            conversations.map((conv) => (
              <Flex
                key={conv.id}
                align="center"
                justify="space-between"
                px={3}
                py={2}
                borderRadius="10px"
                cursor="pointer"
                bg={activeConvId === conv.id ? '#007aff' : 'transparent'}
                color={activeConvId === conv.id ? 'white' : '#1d1d1f'}
                _hover={{ bg: activeConvId === conv.id ? '#007aff' : '#f5f5f7' }}
                transition="all 0.15s"
                onClick={() => setActiveConvId(conv.id)}
              >
                <Text fontSize="13px" fontWeight={activeConvId === conv.id ? '600' : '400'} truncate flex={1}>
                  {conv.title}
                </Text>
                <Box
                  as="button"
                  p={1}
                  borderRadius="6px"
                  opacity={0.5}
                  _hover={{ opacity: 1, color: activeConvId === conv.id ? 'white' : '#ff3b30' }}
                  onClick={(e) => handleDelete(conv.id, e)}
                  flexShrink={0}
                  ml={1}
                >
                  <LuTrash2 size={13} />
                </Box>
              </Flex>
            ))
          ) : (
            <Flex justify="center" py={8}>
              <Text fontSize="12px" color="#86868b">Nessuna conversazione</Text>
            </Flex>
          )}
        </Flex>
      </Flex>

      {/* Main Chat Area */}
      <Flex
        direction="column"
        flex={1}
        bg="white"
        borderRadius="16px"
        boxShadow="0 2px 20px rgba(0,0,0,0.06)"
        overflow="hidden"
      >
        {activeConvId ? (
          <>
            {/* Messages */}
            <Flex direction="column" flex={1} overflow="auto" px={6} py={4}>
              {messages?.map((msg) => (
                <MessageBubble key={msg.id} message={msg} />
              ))}
              {streaming && !streamContent && <TypingIndicator />}
              {streaming && <StreamingBubble content={streamContent} />}
              <div ref={messagesEndRef} />
            </Flex>

            {/* Input */}
            <Box px={4} pb={4} pt={2}>
              <Flex
                bg="#f5f5f7"
                borderRadius="14px"
                align="center"
                px={4}
                gap={2}
              >
                <Input
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && handleSend()}
                  placeholder="Scrivi un messaggio..."
                  bg="transparent"
                  border="none"
                  h="48px"
                  fontSize="14px"
                  _focus={{ boxShadow: 'none' }}
                  flex={1}
                  disabled={streaming}
                />
                <Box
                  as="button"
                  p={2}
                  borderRadius="10px"
                  bg={inputValue.trim() ? '#007aff' : 'transparent'}
                  color={inputValue.trim() ? 'white' : '#86868b'}
                  cursor={inputValue.trim() ? 'pointer' : 'default'}
                  transition="all 0.15s"
                  onClick={handleSend}
                >
                  <LuSend size={18} />
                </Box>
              </Flex>
            </Box>
          </>
        ) : (
          <Flex direction="column" align="center" justify="center" flex={1} gap={3}>
            <LuMessageCircle size={48} color="#007aff" />
            <Text fontSize="17px" fontWeight="600" color="#1d1d1f">Chat AI</Text>
            <Text fontSize="14px" color="#86868b" textAlign="center" maxW="400px">
              Interroga i processi e i documenti della tua azienda.
              Seleziona una conversazione o creane una nuova.
            </Text>
            <Button
              mt={2}
              bg="#007aff"
              color="white"
              borderRadius="10px"
              h="36px"
              px={5}
              fontSize="13px"
              fontWeight="600"
              _hover={{ bg: '#0066d6' }}
              onClick={handleNewChat}
            >
              <Flex align="center" gap={1.5}><LuPlus size={15} /> Nuova Chat</Flex>
            </Button>
          </Flex>
        )}
      </Flex>
    </Flex>
  )
}
