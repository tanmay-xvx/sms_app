/** @type {import('next').NextConfig} */
const nextConfig = {

  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    NEXT_PUBLIC_AI_SERVICE_URL: process.env.NEXT_PUBLIC_AI_SERVICE_URL,
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.NEXT_PUBLIC_API_URL}/api/:path*`,
      },
      {
        source: '/ai/:path*',
        destination: `${process.env.NEXT_PUBLIC_AI_SERVICE_URL}/:path*`,
      },
    ];
  },
};

module.exports = nextConfig; 