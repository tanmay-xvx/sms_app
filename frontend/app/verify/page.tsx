'use client';

import React, { useState, useEffect } from 'react';
import { Shield, ArrowLeft, CheckCircle } from 'lucide-react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import Button from '../components/Button';
import Input from '../components/Input';
import Card from '../components/Card';

export default function VerifyPage() {
  const [otp, setOtp] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [isVerified, setIsVerified] = useState(false);
  const router = useRouter();

  useEffect(() => {
    // Get phone number from localStorage
    const storedPhone = localStorage.getItem('phoneNumber');
    if (!storedPhone) {
      toast.error('Phone number not found. Please start over.');
      router.push('/');
      return;
    }
    setPhoneNumber(storedPhone);
  }, [router]);

  const handleVerifyOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!otp.trim()) {
      setError('OTP is required');
      return;
    }

    if (otp.length !== 6) {
      setError('OTP must be 6 digits');
      return;
    }

    setIsLoading(true);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/sms/verify-otp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phone_number: phoneNumber,
          otp: otp.trim(),
        }),
      });

      const data = await response.json();

              if (response.ok && data.success) {
          toast.success('OTP verified successfully!');
          setIsVerified(true);
          // Clear phone number from localStorage
          localStorage.removeItem('phoneNumber');
          // Set verification flag
          localStorage.setItem('isVerified', 'true');
          sessionStorage.setItem('isVerified', 'true');
          
          // Redirect to dashboard after a short delay
          setTimeout(() => {
            router.push('/dashboard');
          }, 1500);
        } else {
        toast.error(data.message || 'Invalid OTP');
        setError(data.message || 'Invalid OTP');
      }
    } catch (err) {
      console.error('Error verifying OTP:', err);
      toast.error('Network error. Please try again.');
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    setIsLoading(true);
    setError('');

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/sms/send-otp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phone_number: phoneNumber,
        }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        toast.success('New OTP sent successfully!');
        setOtp('');
      } else {
        toast.error(data.message || 'Failed to resend OTP');
        setError(data.message || 'Failed to resend OTP');
      }
    } catch (err) {
      console.error('Error resending OTP:', err);
      toast.error('Network error. Please try again.');
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  if (isVerified) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <Card className="text-center">
            <div className="mb-6">
              <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
                <CheckCircle className="w-8 h-8 text-green-600" />
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Verification Successful!</h1>
              <p className="text-gray-600">Redirecting to dashboard...</p>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div className="bg-green-600 h-2 rounded-full animate-pulse"></div>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <Card title="Verify OTP" className="text-center">
          <div className="mb-6">
            <div className="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-4">
              <Shield className="w-8 h-8 text-blue-600" />
            </div>
            <p className="text-gray-600 mb-4">
              Enter the 6-digit OTP sent to
            </p>
            <p className="text-lg font-semibold text-gray-900 mb-6">
              {phoneNumber}
            </p>
          </div>

          <form onSubmit={handleVerifyOTP} className="space-y-6">
            <Input
              label="OTP Code"
              type="text"
              placeholder="123456"
              value={otp}
              onChange={setOtp}
              error={error}
              required
              maxLength={6}
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
                  <span>Verifying...</span>
                </>
              ) : (
                <>
                  <Shield className="w-4 h-4" />
                  <span>Verify OTP</span>
                </>
              )}
            </Button>
          </form>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-500 mb-4">
              Didn't receive the OTP?
            </p>
            <Button
              variant="outline"
              size="sm"
              onClick={handleResendOTP}
              disabled={isLoading}
              className="w-full"
            >
              Resend OTP
            </Button>
          </div>

          <div className="mt-4">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => router.push('/')}
              className="w-full flex items-center justify-center space-x-2"
            >
              <ArrowLeft className="w-4 h-4" />
              <span>Back to Start</span>
            </Button>
          </div>
        </Card>
      </div>
    </div>
  );
} 