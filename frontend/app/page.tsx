'use client';

import React, { useState } from 'react';
import { Phone, Send, ArrowRight } from 'lucide-react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import Button from './components/Button';
import Input from './components/Input';
import Card from './components/Card';

export default function HomePage() {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const validatePhoneNumber = (phone: string) => {
    // Basic phone validation - should start with + and have 10+ digits
    const phoneRegex = /^\+[1-9]\d{1,14}$/;
    return phoneRegex.test(phone);
  };

  const handleSendOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!phoneNumber.trim()) {
      setError('Phone number is required');
      return;
    }

    if (!validatePhoneNumber(phoneNumber)) {
      setError('Please enter a valid phone number (e.g., +1234567890)');
      return;
    }

    setIsLoading(true);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/sms/send-otp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phone_number: phoneNumber.trim(),
        }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        toast.success('OTP sent successfully!');
        // Store phone number in localStorage for verification page
        localStorage.setItem('phoneNumber', phoneNumber);
        router.push('/verify');
      } else {
        toast.error(data.message || 'Failed to send OTP');
        setError(data.message || 'Failed to send OTP');
      }
    } catch (err) {
      console.error('Error sending OTP:', err);
      toast.error('Network error. Please try again.');
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <Card title="Welcome to SMS App" className="text-center">
          <div className="mb-6">
            <div className="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-4">
              <Phone className="w-8 h-8 text-blue-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Get Started</h1>
            <p className="text-gray-600">Enter your phone number to receive an OTP</p>
          </div>

          <form onSubmit={handleSendOTP} className="space-y-6">
            <Input
              label="Phone Number"
              type="tel"
              placeholder="+1234567890"
              value={phoneNumber}
              onChange={setPhoneNumber}
              error={error}
              required
            />

            <Button
              type="submit"
              disabled={isLoading}
              className="w-full flex items-center justify-center space-x-2"
              size="lg"
            >
              {isLoading ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <span>Sending...</span>
                </>
              ) : (
                <>
                  <Send className="w-4 h-4" />
                  <span>Send OTP</span>
                </>
              )}
            </Button>
          </form>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-500">
              By continuing, you agree to our terms of service
            </p>
          </div>
        </Card>

        {/* Demo section */}
        <div className="mt-8 text-center">
          <p className="text-sm text-gray-500 mb-2">Demo Features</p>
          <div className="flex justify-center space-x-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push('/dashboard')}
              className="flex items-center space-x-1"
            >
              <span>Dashboard</span>
              <ArrowRight className="w-3 h-3" />
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push('/admin')}
              className="flex items-center space-x-1"
            >
              <span>Admin</span>
              <ArrowRight className="w-3 h-3" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
} 