import { api } from '@/Core/Services/Api'
import type { Conversation, Message } from '../Types'

const prefix = 'chat'

const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    conversations: builder.query<Conversation[], void>({
      query: () => `${prefix}/conversations`,
      providesTags: ['Conversations', 'HasCompanyHeader'],
    }),
    conversationMessages: builder.query<Message[], number>({
      query: (id) => `${prefix}/conversations/${id}/messages`,
      providesTags: (_r, _e, id) => [{ type: 'Messages' as const, id }],
    }),
    createConversation: builder.mutation<Conversation, { title?: string }>({
      query: (body) => ({
        url: `${prefix}/conversations`,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Conversations'],
    }),
    deleteConversation: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/conversations/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Conversations'],
    }),
  }),
  overrideExisting: false,
})

export const {
  useConversationsQuery,
  useConversationMessagesQuery,
  useCreateConversationMutation,
  useDeleteConversationMutation,
} = extendedApi
