import { Box, Button, Input, Text, VStack } from '@chakra-ui/react'
import { useState } from 'react'
import { useDispatch } from 'react-redux'
import { useNavigate } from 'react-router-dom'

import Config from '@/Config'
import Logger from '@/Core/Services/Logger'
import { setToken, setUser } from '../Redux'
import { useSignInMutation, useLazyCurrentUserQuery } from '../Services/Api'

export default function SignInView() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [signIn, { isLoading }] = useSignInMutation()
  const [getCurrentUser] = useLazyCurrentUserQuery()
  const dispatch = useDispatch()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    try {
      const result = await signIn({ email, password }).unwrap()
      dispatch(setToken(result))
      const user = await getCurrentUser().unwrap()
      dispatch(setUser(user))
      navigate(Config.urls.home)
    } catch (err) {
      Logger.error('Sign in error', err)
      setError('Credenziali non valide')
    }
  }

  return (
    <Box display="flex" alignItems="center" justifyContent="center" minH="100vh" bg="#f5f5f7">
      <Box
        as="form"
        onSubmit={handleSubmit}
        w="full"
        maxW="380px"
        p={8}
        bg="white"
        borderRadius="16px"
        boxShadow="0 2px 16px rgba(0,0,0,0.08)"
      >
        <VStack gap={5}>
          <Text fontSize="24px" fontWeight="700" color="#1d1d1f">Orbita</Text>
          <Text fontSize="13px" color="#86868b">Accedi con le tue credenziali</Text>

          <Input
            name="email"
            autoComplete="username"
            placeholder="Email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            bg="#f5f5f7"
            border="none"
            borderRadius="10px"
            h="44px"
            fontSize="14px"
            _focus={{ bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.3)' }}
          />
          <Input
            name="password"
            autoComplete="current-password"
            placeholder="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            bg="#f5f5f7"
            border="none"
            borderRadius="10px"
            h="44px"
            fontSize="14px"
            _focus={{ bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.3)' }}
          />

          {error && <Text color="#ff3b30" fontSize="13px">{error}</Text>}

          <Button
            type="submit"
            w="full"
            h="44px"
            bg="#007aff"
            color="white"
            borderRadius="10px"
            fontSize="15px"
            fontWeight="600"
            loading={isLoading}
            _hover={{ bg: '#0066d6' }}
          >
            Accedi
          </Button>
        </VStack>
      </Box>
    </Box>
  )
}
