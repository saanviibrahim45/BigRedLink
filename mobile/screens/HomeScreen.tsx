import React from 'react'
import { View, Text, Button, StyleSheet } from 'react-native'
import AuthService from '../services/AuthService'

interface HomeScreenProps {
  navigation: {
    replace: (screen: 'Login' | 'Home') => void
  }
}

const HomeScreen: React.FC<HomeScreenProps> = ({ navigation }) => {
  const handleLogout = async () => {
    try {
      await AuthService.logout()
      navigation.replace('Login')
    } catch (err) {
      console.error('Logout failed', err)
    }
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Welcome to Big Red Link!</Text>
      <Button title="Logout" onPress={handleLogout} />
    </View>
  )
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  title: { fontSize: 20, fontWeight: 'bold', marginBottom: 20 },
})

export default HomeScreen
