import { Box, Button, Flex, Input, Text } from '@chakra-ui/react'
import { useState } from 'react'

import { useCurrentUser } from '@/Auth/Hooks'
import { useUpdateUserPasswordMutation } from '../Services/Api'

export default function PasswordView() {
  const user = useCurrentUser()
  const [updatePassword, { isLoading }] = useUpdateUserPasswordMutation()
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setMessage('')

    if (newPassword !== confirmPassword) {
      setError('Le password non coincidono')
      return
    }
    if (newPassword.length < 8) {
      setError('La password deve essere di almeno 8 caratteri')
      return
    }

    try {
      await updatePassword({
        id: user!.id,
        body: { currentPassword, password: newPassword },
      }).unwrap()
      setMessage('Password aggiornata con successo')
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch {
      setError('Errore durante l\'aggiornamento della password')
    }
  }

  const inputStyle = {
    bg: '#f5f5f7',
    border: 'none',
    borderRadius: '10px',
    h: '44px',
    fontSize: '14px',
    _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
  }

  return (
    <Box>
      <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Reimposta Password</Text>

      <Box
        as="form"
        onSubmit={handleSubmit}
        maxW="400px"
      >
        <Flex direction="column" gap={4}>
          <Box>
            <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1.5}>Password attuale</Text>
            <Input
              type="password"
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              required
              {...inputStyle}
            />
          </Box>
          <Box>
            <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1.5}>Nuova password</Text>
            <Input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              {...inputStyle}
            />
          </Box>
          <Box>
            <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1.5}>Conferma nuova password</Text>
            <Input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              {...inputStyle}
            />
          </Box>

          {error && <Text fontSize="13px" color="#ff3b30">{error}</Text>}
          {message && <Text fontSize="13px" color="#34c759">{message}</Text>}

          <Button
            type="submit"
            bg="#007aff"
            color="white"
            borderRadius="10px"
            h="44px"
            fontSize="14px"
            fontWeight="600"
            loading={isLoading}
            _hover={{ bg: '#0066d6' }}
          >
            Aggiorna Password
          </Button>
        </Flex>
      </Box>
    </Box>
  )
}
