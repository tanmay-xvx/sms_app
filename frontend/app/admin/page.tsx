'use client';

import React, { useState, useEffect } from 'react';
import { BarChart3, RefreshCw, Download, Eye, Phone, MessageSquare, Clock, Shield } from 'lucide-react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import Button from '../components/Button';
import Card from '../components/Card';

interface LogData {
  otps: {
    count: number;
    data: any[];
  };
  callbacks: {
    count: number;
    data: any[];
  };
  sms: {
    count: number;
    data: any[];
  };
  timestamp: string;
  total_records: number;
}

export default function AdminPage() {
  const [logs, setLogs] = useState<LogData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [limit, setLimit] = useState(100);
  const [activeTab, setActiveTab] = useState<'overview' | 'otps' | 'callbacks' | 'sms'>('overview');
  const router = useRouter();

  useEffect(() => {
    fetchLogs();
  }, [limit]);

  const fetchLogs = async () => {
    setIsLoading(true);
    setError('');

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/logs?limit=${limit}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch logs');
      }

      const data = await response.json();
      setLogs(data);
    } catch (err) {
      console.error('Error fetching logs:', err);
      setError('Failed to fetch logs. Please try again.');
      toast.error('Failed to fetch logs');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'sent':
      case 'delivered':
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'pending':
      case 'requested':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const exportLogs = () => {
    if (!logs) return;
    
    const dataStr = JSON.stringify(logs, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `sms-logs-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
    URL.revokeObjectURL(url);
    toast.success('Logs exported successfully');
  };

  if (isLoading && !logs) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">Loading logs...</p>
        </div>
      </div>
    );
  }

  if (error && !logs) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">{error}</p>
          <Button onClick={fetchLogs}>Retry</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <BarChart3 className="w-8 h-8 text-blue-600" />
              <h1 className="text-xl font-semibold text-gray-900">Admin Dashboard</h1>
            </div>
            <div className="flex items-center space-x-3">
              <select
                value={limit}
                onChange={(e) => setLimit(Number(e.target.value))}
                className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value={50}>50 records</option>
                <option value={100}>100 records</option>
                <option value={200}>200 records</option>
              </select>
              <Button
                variant="outline"
                size="sm"
                onClick={fetchLogs}
                disabled={isLoading}
                className="flex items-center space-x-2"
              >
                <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
                <span>Refresh</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={exportLogs}
                className="flex items-center space-x-2"
              >
                <Download className="w-4 h-4" />
                <span>Export</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Overview Cards */}
        <div className="grid md:grid-cols-4 gap-6 mb-8">
          <Card>
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                <MessageSquare className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Total SMS</p>
                <p className="text-2xl font-bold text-gray-900">{logs?.sms.count || 0}</p>
              </div>
            </div>
          </Card>

          <Card>
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                <Phone className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Callbacks</p>
                <p className="text-2xl font-bold text-gray-900">{logs?.callbacks.count || 0}</p>
              </div>
            </div>
          </Card>

          <Card>
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
                <Shield className="w-5 h-5 text-yellow-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Active OTPs</p>
                <p className="text-2xl font-bold text-gray-900">{logs?.otps.count || 0}</p>
              </div>
            </div>
          </Card>

          <Card>
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                <Clock className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Last Updated</p>
                <p className="text-sm font-semibold text-gray-900">
                  {logs?.timestamp ? formatTimestamp(logs.timestamp) : 'N/A'}
                </p>
              </div>
            </div>
          </Card>
        </div>

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="border-b border-gray-200">
            <nav className="flex space-x-8 px-6">
              {[
                { id: 'overview', label: 'Overview', count: logs?.total_records || 0 },
                { id: 'sms', label: 'SMS Logs', count: logs?.sms.count || 0 },
                { id: 'callbacks', label: 'Callbacks', count: logs?.callbacks.count || 0 },
                { id: 'otps', label: 'OTPs', count: logs?.otps.count || 0 },
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`py-4 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  {tab.label}
                  <span className="ml-2 bg-gray-100 text-gray-900 py-0.5 px-2.5 rounded-full text-xs font-medium">
                    {tab.count}
                  </span>
                </button>
              ))}
            </nav>
          </div>

          <div className="p-6">
            {/* Overview Tab */}
            {activeTab === 'overview' && (
              <div className="space-y-6">
                <h3 className="text-lg font-medium text-gray-900">System Overview</h3>
                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <h4 className="font-medium text-gray-700 mb-3">Recent Activity</h4>
                    <div className="space-y-2">
                      <p className="text-sm text-gray-600">
                        <span className="font-medium">Total Records:</span> {logs?.total_records || 0}
                      </p>
                      <p className="text-sm text-gray-600">
                        <span className="font-medium">Last Updated:</span> {logs?.timestamp ? formatTimestamp(logs.timestamp) : 'N/A'}
                      </p>
                    </div>
                  </div>
                  <div>
                    <h4 className="font-medium text-gray-700 mb-3">Quick Actions</h4>
                    <div className="space-y-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => router.push('/dashboard')}
                        className="w-full justify-center"
                      >
                        Go to Dashboard
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => router.push('/')}
                        className="w-full justify-center"
                      >
                        Send OTP
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* SMS Tab */}
            {activeTab === 'sms' && (
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">SMS Logs</h3>
                {logs?.sms.data && logs.sms.data.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">To</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Message</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Sent At</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {logs.sms.data.map((sms: any) => (
                          <tr key={sms.id} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {sms.id.slice(-8)}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{sms.from}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{sms.to}</td>
                            <td className="px-6 py-4 text-sm text-gray-900 max-w-xs truncate">{sms.message}</td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(sms.status)}`}>
                                {sms.status}
                              </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {formatTimestamp(sms.sent_at)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <p className="text-gray-500 text-center py-8">No SMS logs found</p>
                )}
              </div>
            )}

            {/* Callbacks Tab */}
            {activeTab === 'callbacks' && (
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">Callback Requests</h3>
                {logs?.callbacks.data && logs.callbacks.data.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Phone</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Message</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Priority</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Requested At</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {logs.callbacks.data.map((callback: any) => (
                          <tr key={callback.id} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {callback.id.slice(-8)}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{callback.phone_number}</td>
                            <td className="px-6 py-4 text-sm text-gray-900 max-w-xs truncate">{callback.message}</td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                                callback.priority === 'urgent' ? 'bg-red-100 text-red-800' :
                                callback.priority === 'high' ? 'bg-orange-100 text-orange-800' :
                                'bg-blue-100 text-blue-800'
                              }`}>
                                {callback.priority}
                              </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(callback.status)}`}>
                                {callback.status}
                              </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {formatTimestamp(callback.requested_at)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <p className="text-gray-500 text-center py-8">No callback requests found</p>
                )}
              </div>
            )}

            {/* OTPs Tab */}
            {activeTab === 'otps' && (
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">OTP Records</h3>
                {logs?.otps.data && logs.otps.data.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Phone</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Attempts</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Expires At</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {logs.otps.data.map((otp: any) => (
                          <tr key={otp.phone} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{otp.phone}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{otp.attempts}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {otp.expires_at ? formatTimestamp(otp.expires_at) : 'N/A'}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {formatTimestamp(otp.created_at)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <p className="text-gray-500 text-center py-8">No OTP records found</p>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
} 