import { defineSemanticTokens } from '@chakra-ui/react'

export const colors = defineSemanticTokens.colors({
  brand: {
    contrast: { value: { _light: 'black', _dark: 'black' } },
    bg: { value: { _light: '{colors.brand.300}', _dark: '{colors.brand.950}' } },
    fg: { value: { _light: '{colors.brand.500}', _dark: '{colors.brand.300}' } },
    subtle: { value: { _light: '{colors.brand.100}', _dark: '{colors.brand.900}' } },
    muted: { value: { _light: '{colors.brand.200}', _dark: '{colors.brand.800}' } },
    emphasized: { value: { _light: '{colors.brand.300}', _dark: '{colors.brand.700}' } },
    solid: { value: { _light: '{colors.brand.300}', _dark: '{colors.brand.300}' } },
    focusRing: { value: { _light: '{colors.brand.500}', _dark: '{colors.brand.500}' } },
  },
  secondary: {
    contrast: { value: { _light: 'black', _dark: '{colors.secondary.100}' } },
    bg: { value: { _light: '{colors.secondary.50}', _dark: '{colors.secondary.950}' } },
    fg: { value: { _light: '{colors.secondary.800}', _dark: '{colors.secondary.300}' } },
    subtle: { value: { _light: '{colors.secondary.100}', _dark: '{colors.secondary.900}' } },
    muted: { value: { _light: '{colors.secondary.200}', _dark: '{colors.secondary.800}' } },
    emphasized: { value: { _light: '{colors.secondary.300}', _dark: '{colors.secondary.700}' } },
    solid: { value: { _light: '{colors.secondary.300}', _dark: '{colors.secondary.300}' } },
    focusRing: { value: { _light: '{colors.secondary.500}', _dark: '{colors.secondary.500}' } },
  },
  border: {
    input: { value: { _light: '{colors.gray.300}', _dark: '{colors.gray.500}' } },
  },
  danger: {
    value: { base: '{colors.red.100}', _dark: '{colors.red.700}' },
  },
  success: {
    value: { base: '{colors.green.100}', _dark: '{colors.green.700}' },
  },
})
