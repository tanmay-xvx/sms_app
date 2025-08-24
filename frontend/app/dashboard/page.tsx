'use client';

import React, { useState, useEffect } from 'react';
import { Phone, MessageSquare, Send, User, LogOut } from 'lucide-react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import Button from '../components/Button';
import Input from '../components/Input';
import Card from '../components/Card';

interface ChatMessage {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
}

export default function DashboardPage() {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [message, setMessage] = useState('');
  const [priority, setPriority] = useState('normal');
  const [isLoadingCallback, setIsLoadingCallback] = useState(false);
  const [isLoadingChat, setIsLoadingChat] = useState(false);
  const [chatInput, setChatInput] = useState('');
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [showCallbackForm, setShowCallbackForm] = useState(false);
  const router = useRouter();

  useEffect(() => {
    // Check if user is verified (in a real app, you'd check JWT token)
    const isVerified = localStorage.getItem('isVerified') || sessionStorage.getItem('isVerified');
    if (!isVerified) {
      toast.error('Please verify your OTP first');
      router.push('/');
      return;
    }

    // Add welcome message
    setChatMessages([
      {
        id: '1',
        type: 'ai',
        content: 'Hello! I\'m your AI assistant. How can I help you today?',
        timestamp: new Date(),
      },
    ]);
  }, [router]);

  const handleRequestCallback = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoadingCallback(true);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/callback/request`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phone_number: phoneNumber,
          message: message,
          priority: priority,
        }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        toast.success('Callback requested successfully!');
        setShowCallbackForm(false);
        setPhoneNumber('');
        setMessage('');
        setPriority('normal');
        
        // Add success message to chat
        setChatMessages(prev => [...prev, {
          id: Date.now().toString(),
          type: 'ai',
          content: `Great! I've requested a callback for you. Request ID: ${data.request_id}. Our team will contact you soon.`,
          timestamp: new Date(),
        }]);
      } else {
        toast.error(data.message || 'Failed to request callback');
      }
    } catch (err) {
      console.error('Error requesting callback:', err);
      toast.error('Network error. Please try again.');
    } finally {
      setIsLoadingCallback(false);
    }
  };

  const handleSendChatMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!chatInput.trim()) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: chatInput.trim(),
      timestamp: new Date(),
    };

    setChatMessages(prev => [...prev, userMessage]);
    setChatInput('');
    setIsLoadingChat(true);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_AI_SERVICE_URL || 'http://localhost:8000'}/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          question: userMessage.content,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        const aiMessage: ChatMessage = {
          id: (Date.now() + 1).toString(),
          type: 'ai',
          content: data.answer,
          timestamp: new Date(),
        };

        setChatMessages(prev => [...prev, aiMessage]);

        // If escalation is needed, suggest callback
        if (data.escalate) {
          setTimeout(() => {
            setChatMessages(prev => [...prev, {
              id: (Date.now() + 2).toString(),
              type: 'ai',
              content: 'It looks like you might need human assistance. Would you like me to request a callback from our support team?',
              timestamp: new Date(),
            }]);
          }, 1000);
        }
      } else {
        toast.error('Failed to get AI response');
      }
    } catch (err) {
      console.error('Error getting AI response:', err);
      toast.error('Network error. Please try again.');
      
      // Add error message to chat
      setChatMessages(prev => [...prev, {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: 'Sorry, I\'m having trouble connecting right now. Please try again later.',
        timestamp: new Date(),
      }]);
    } finally {
      setIsLoadingChat(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('isVerified');
    sessionStorage.removeItem('isVerified');
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <MessageSquare className="w-8 h-8 text-blue-600" />
              <h1 className="text-xl font-semibold text-gray-900">SMS Dashboard</h1>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleLogout}
              className="flex items-center space-x-2"
            >
              <LogOut className="w-4 h-4" />
              <span>Logout</span>
            </Button>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid lg:grid-cols-2 gap-8">
          {/* Left Column - Callback Request */}
          <div className="space-y-6">
            <Card title="Request Callback">
              <div className="space-y-4">
                <p className="text-gray-600">
                  Need human assistance? Request a callback from our support team.
                </p>
                
                {!showCallbackForm ? (
                  <Button
                    onClick={() => setShowCallbackForm(true)}
                    className="w-full flex items-center justify-center space-x-2"
                  >
                    <Phone className="w-4 h-4" />
                    <span>Request Callback</span>
                  </Button>
                ) : (
                  <form onSubmit={handleRequestCallback} className="space-y-4">
                    <Input
                      label="Phone Number"
                      type="tel"
                      placeholder="+1234567890"
                      value={phoneNumber}
                      onChange={setPhoneNumber}
                      required
                    />
                    
                    <Input
                      label="Message"
                      type="text"
                      placeholder="Brief description of your issue"
                      value={message}
                      onChange={setMessage}
                      required
                    />
                    
                    <div className="space-y-2">
                      <label className="block text-sm font-medium text-gray-700">
                        Priority
                      </label>
                      <select
                        value={priority}
                        onChange={(e) => setPriority(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      >
                        <option value="normal">Normal</option>
                        <option value="high">High</option>
                        <option value="urgent">Urgent</option>
                      </select>
                    </div>
                    
                    <div className="flex space-x-3">
                      <Button
                        type="submit"
                        disabled={isLoadingCallback}
                        className="flex-1"
                      >
                        {isLoadingCallback ? 'Requesting...' : 'Request Callback'}
                      </Button>
                      <Button
                        type="button"
                        variant="secondary"
                        onClick={() => setShowCallbackForm(false)}
                        className="flex-1"
                      >
                        Cancel
                      </Button>
                    </div>
                  </form>
                )}
              </div>
            </Card>

            <Card title="Quick Actions">
              <div className="grid grid-cols-2 gap-3">
                <Button
                  variant="outline"
                  onClick={() => router.push('/admin')}
                  className="flex items-center justify-center space-x-2"
                >
                  <User className="w-4 h-4" />
                  <span>View Logs</span>
                </Button>
                <Button
                  variant="outline"
                  onClick={() => router.push('/')}
                  className="flex items-center justify-center space-x-2"
                >
                  <Phone className="w-4 h-4" />
                  <span>Send OTP</span>
                </Button>
              </div>
            </Card>
          </div>

          {/* Right Column - Chat */}
          <div className="space-y-6">
            <Card title="AI Assistant">
              <div className="space-y-4">
                {/* Chat Messages */}
                <div className="h-96 overflow-y-auto space-y-3 p-3 bg-gray-50 rounded-lg">
                  {chatMessages.map((msg) => (
                    <div
                      key={msg.id}
                      className={`flex ${msg.type === 'user' ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={`max-w-xs lg:max-w-md px-3 py-2 rounded-lg ${
                          msg.type === 'user'
                            ? 'bg-blue-600 text-white'
                            : 'bg-white text-gray-800 border border-gray-200'
                        }`}
                      >
                        <p className="text-sm">{msg.content}</p>
                        <p className={`text-xs mt-1 ${
                          msg.type === 'user' ? 'text-blue-100' : 'text-gray-500'
                        }`}>
                          {msg.timestamp.toLocaleTimeString()}
                        </p>
                      </div>
                    </div>
                  ))}
                  {isLoadingChat && (
                    <div className="flex justify-start">
                      <div className="bg-white text-gray-800 border border-gray-200 px-3 py-2 rounded-lg">
                        <div className="flex space-x-1">
                          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                        </div>
                      </div>
                    </div>
                  )}
                </div>

                {/* Chat Input */}
                <form onSubmit={handleSendChatMessage} className="flex space-x-2">
                  <Input
                    placeholder="Type your question..."
                    value={chatInput}
                    onChange={setChatInput}
                    className="flex-1"
                  />
                  <Button
                    type="submit"
                    disabled={isLoadingChat || !chatInput.trim()}
                    size="sm"
                  >
                    <Send className="w-4 h-4" />
                  </Button>
                </form>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
} 