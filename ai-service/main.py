from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing import List, Optional
import os
import httpx
from dotenv import load_dotenv
import asyncio
import logging

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="SMS AI Service",
    description="AI microservice for SMS application",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# CORS middleware
cors_origins = ["http://localhost:3000", "http://localhost:8080"]
if os.getenv("ENVIRONMENT") == "production":
    if additional_origins := os.getenv("CORS_ORIGINS"):
        cors_origins.extend(additional_origins.split(","))

app.add_middleware(
    CORSMiddleware,
    allow_origins=cors_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic models
class MessageAnalysisRequest(BaseModel):
    message: str
    user_id: Optional[str] = None

class MessageAnalysisResponse(BaseModel):
    sentiment: str
    intent: str
    confidence: float
    keywords: List[str]
    suggested_response: Optional[str] = None

class SummarizationRequest(BaseModel):
    messages: List[str]
    max_length: Optional[int] = 100

class SummarizationResponse(BaseModel):
    summary: str
    key_points: List[str]
    word_count: int

class ChatRequest(BaseModel):
    question: str

class ChatResponse(BaseModel):
    answer: str
    confidence: float
    escalate: bool = False
    source: str = "ai"
    suggested_actions: Optional[List[str]] = None

class HealthResponse(BaseModel):
    status: str
    service: str
    version: str

# Configuration
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
MODEL_NAME = os.getenv("MODEL_NAME", "gpt-3.5-turbo")

@app.get("/", response_model=HealthResponse)
async def root():
    """Root endpoint with service information"""
    return HealthResponse(
        status="healthy",
        service="sms-ai-service",
        version="1.0.0"
    )

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        service="sms-ai-service",
        version="1.0.0"
    )

@app.post("/analyze", response_model=MessageAnalysisResponse)
async def analyze_message(request: MessageAnalysisRequest):
    """
    Analyze a message for sentiment, intent, and other insights
    """
    try:
        if not OPENAI_API_KEY:
            # Fallback to mock analysis if no API key
            return await mock_message_analysis(request.message)
        
        # TODO: Implement OpenAI API call
        # For now, return mock data
        return await mock_message_analysis(request.message)
        
    except Exception as e:
        logger.error(f"Error analyzing message: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to analyze message")

@app.post("/summarize", response_model=SummarizationResponse)
async def summarize_messages(request: SummarizationRequest):
    """
    Summarize a list of messages
    """
    try:
        if not OPENAI_API_KEY:
            # Fallback to mock summarization if no API key
            return await mock_summarization(request.messages, request.max_length)
        
        # TODO: Implement OpenAI API call
        # For now, return mock data
        return await mock_summarization(request.messages, request.max_length)
        
    except Exception as e:
        logger.error(f"Error summarizing messages: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to summarize messages")

@app.post("/generate-response")
async def generate_response(request: MessageAnalysisRequest):
    """
    Generate an AI-powered response to a message
    """
    try:
        if not OPENAI_API_KEY:
            # Fallback to mock response generation
            return {"response": f"Thank you for your message: '{request.message}'. How can I help you today?"}
        
        # TODO: Implement OpenAI API call for response generation
        return {"response": f"Thank you for your message: '{request.message}'. How can I help you today?"}
        
    except Exception as e:
        logger.error(f"Error generating response: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to generate response")

@app.post("/classify-intent")
async def classify_intent(request: MessageAnalysisRequest):
    """
    Classify the intent of a message
    """
    try:
        # Simple intent classification logic
        message_lower = request.message.lower()
        
        if any(word in message_lower for word in ["hello", "hi", "hey", "good morning"]):
            intent = "greeting"
        elif any(word in message_lower for word in ["help", "support", "assist"]):
            intent = "help_request"
        elif any(word in message_lower for word in ["complaint", "issue", "problem"]):
            intent = "complaint"
        elif any(word in message_lower for word in ["thank", "thanks", "appreciate"]):
            intent = "gratitude"
        else:
            intent = "general"
        
        return {
            "intent": intent,
            "confidence": 0.85,
            "message": request.message
        }
        
    except Exception as e:
        logger.error(f"Error classifying intent: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to classify intent")

@app.post("/chat", response_model=ChatResponse)
async def chat_with_ai(request: ChatRequest):
    """
    Chat with AI to get FAQ-like answers
    Returns escalation flag if confidence is low
    """
    try:
        if not OPENAI_API_KEY:
            # Fallback to mock chat if no API key
            return await mock_chat_response(request.question)
        
        # Use OpenAI API for real responses
        return await openai_chat_response(request.question)
        
    except Exception as e:
        logger.error(f"Error in chat endpoint: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to process chat request")

# Mock functions for development/testing
async def mock_message_analysis(message: str) -> MessageAnalysisResponse:
    """Mock message analysis for development"""
    await asyncio.sleep(0.1)  # Simulate processing time
    
    # Simple sentiment analysis
    positive_words = ["good", "great", "excellent", "amazing", "love", "happy"]
    negative_words = ["bad", "terrible", "awful", "hate", "sad", "angry"]
    
    message_lower = message.lower()
    positive_count = sum(1 for word in positive_words if word in message_lower)
    negative_count = sum(1 for word in negative_words if word in message_lower)
    
    if positive_count > negative_count:
        sentiment = "positive"
        confidence = 0.8
    elif negative_count > positive_count:
        sentiment = "negative"
        confidence = 0.8
    else:
        sentiment = "neutral"
        confidence = 0.6
    
    # Simple intent detection
    if any(word in message_lower for word in ["hello", "hi", "hey"]):
        intent = "greeting"
    elif any(word in message_lower for word in ["help", "support"]):
        intent = "help_request"
    else:
        intent = "general"
    
    # Extract keywords (simple approach)
    keywords = [word for word in message_lower.split() if len(word) > 3]
    
    return MessageAnalysisResponse(
        sentiment=sentiment,
        intent=intent,
        confidence=confidence,
        keywords=keywords[:5],  # Limit to 5 keywords
        suggested_response=f"Thank you for your {sentiment} message about {intent}."
    )

async def mock_summarization(messages: List[str], max_length: Optional[int] = 100) -> SummarizationResponse:
    """Mock summarization for development"""
    await asyncio.sleep(0.2)  # Simulate processing time
    
    if not messages:
        return SummarizationResponse(
            summary="No messages to summarize",
            key_points=[],
            word_count=0
        )
    
    # Simple summarization logic
    all_text = " ".join(messages)
    words = all_text.split()
    
    # Create a simple summary
    if len(words) <= max_length:
        summary = all_text
    else:
        summary = " ".join(words[:max_length]) + "..."
    
    # Extract key points (simple approach)
    key_points = []
    for i, msg in enumerate(messages[:3]):  # Limit to first 3 messages
        if len(msg) > 20:
            key_points.append(f"Message {i+1}: {msg[:50]}...")
        else:
            key_points.append(f"Message {i+1}: {msg}")
    
    return SummarizationResponse(
        summary=summary,
        key_points=key_points,
        word_count=len(words)
    )

async def openai_chat_response(question: str) -> ChatResponse:
    """Generate response using OpenAI API"""
    try:
        import openai
        
        # Configure OpenAI client
        client = openai.OpenAI(api_key=OPENAI_API_KEY)
        
        # Create system prompt for FAQ-like responses
        system_prompt = """You are a helpful customer service AI for an SMS application. 
        Answer questions in a friendly, helpful manner. If you're not confident about an answer 
        or if the question requires human intervention, indicate that escalation is needed.
        
        Common topics you can help with:
        - SMS features and usage
        - Account management
        - Technical support basics
        - General inquiries
        
        If you're not sure about something or if it requires account-specific information, 
        suggest escalating to a human agent."""
        
        # Generate response
        response = client.chat.completions.create(
            model=MODEL_NAME,
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": question}
            ],
            max_tokens=300,
            temperature=0.7
        )
        
        answer = response.choices[0].message.content.strip()
        
        # Analyze confidence based on response characteristics
        confidence = analyze_response_confidence(answer, question)
        
        # Determine if escalation is needed
        escalate = confidence < 0.6 or should_escalate(question, answer)
        
        suggested_actions = []
        if escalate:
            suggested_actions = [
                "Contact customer support",
                "Request callback from agent",
                "Check account status"
            ]
        
        return ChatResponse(
            answer=answer,
            confidence=confidence,
            escalate=escalate,
            source="openai",
            suggested_actions=suggested_actions
        )
        
    except Exception as e:
        logger.error(f"OpenAI API error: {str(e)}")
        # Fallback to mock response on API failure
        return await mock_chat_response(question)

async def mock_chat_response(question: str) -> ChatResponse:
    """Mock chat response for development/testing"""
    await asyncio.sleep(0.2)  # Simulate processing time
    
    question_lower = question.lower()
    
    # Simple FAQ responses
    faq_responses = {
        "password": {
            "answer": "To reset your password, you can use the 'Forgot Password' link on the login page, or contact our support team for assistance.",
            "confidence": 0.9,
            "escalate": False
        },
        "account": {
            "answer": "You can manage your account settings through the account dashboard. For account-specific issues, please contact support.",
            "confidence": 0.8,
            "escalate": False
        },
        "sms": {
            "answer": "Our SMS service allows you to send messages to any phone number. You can also schedule messages and view delivery status.",
            "confidence": 0.9,
            "escalate": False
        },
        "support": {
            "answer": "For technical support, please check our help center first. If you need immediate assistance, we can connect you with a human agent.",
            "confidence": 0.7,
            "escalate": True
        },
        "billing": {
            "answer": "Billing inquiries require access to your specific account information. Please contact our billing department for assistance.",
            "confidence": 0.5,
            "escalate": True
        },
        "callback": {
            "answer": "I can help you request a callback from our support team. Would you like me to initiate that process for you?",
            "confidence": 0.8,
            "escalate": False
        }
    }
    
    # Find best matching FAQ
    best_match = None
    best_confidence = 0.0
    
    for keyword, response in faq_responses.items():
        if keyword in question_lower:
            if response["confidence"] > best_confidence:
                best_match = response
                best_confidence = response["confidence"]
    
    # Default response if no FAQ match
    if not best_match:
        best_match = {
            "answer": "I understand you're asking about '" + question + "'. While I can help with general information, this specific question might require human assistance. Would you like me to connect you with our support team?",
            "confidence": 0.4,
            "escalate": True
        }
    
    suggested_actions = []
    if best_match["escalate"]:
        suggested_actions = [
            "Request callback from support agent",
            "Contact customer service",
            "Check our help center"
        ]
    
    return ChatResponse(
        answer=best_match["answer"],
        confidence=best_match["confidence"],
        escalate=best_match["escalate"],
        source="mock",
        suggested_actions=suggested_actions
    )

def analyze_response_confidence(answer: str, question: str) -> float:
    """Analyze confidence level of AI response"""
    # Simple confidence analysis based on response characteristics
    confidence = 0.7  # Base confidence
    
    # Increase confidence for longer, detailed answers
    if len(answer.split()) > 20:
        confidence += 0.1
    
    # Increase confidence if answer directly addresses question keywords
    question_words = set(question.lower().split())
    answer_words = set(answer.lower().split())
    common_words = question_words.intersection(answer_words)
    if len(common_words) > 0:
        confidence += 0.1
    
    # Decrease confidence for vague responses
    vague_phrases = ["i'm not sure", "i don't know", "maybe", "possibly", "could be"]
    if any(phrase in answer.lower() for phrase in vague_phrases):
        confidence -= 0.2
    
    # Decrease confidence for responses suggesting escalation
    escalation_phrases = ["contact support", "human agent", "escalate", "support team"]
    if any(phrase in answer.lower() for phrase in escalation_phrases):
        confidence -= 0.1
    
    return min(max(confidence, 0.0), 1.0)

def should_escalate(question: str, answer: str) -> bool:
    """Determine if escalation is needed based on question and answer"""
    # Keywords that typically require human intervention
    escalation_keywords = [
        "account", "billing", "payment", "refund", "complaint", "dispute",
        "personal", "private", "urgent", "emergency", "legal", "policy"
    ]
    
    question_lower = question.lower()
    answer_lower = answer.lower()
    
    # Check if question contains escalation keywords
    if any(keyword in question_lower for keyword in escalation_keywords):
        return True
    
    # Check if answer suggests escalation
    if any(phrase in answer_lower for phrase in ["contact support", "human agent", "escalate"]):
        return True
    
    return False

if __name__ == "__main__":
    import uvicorn
    
    port = int(os.getenv("PORT", 8000))
    host = os.getenv("HOST", "0.0.0.0")
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=True,
        log_level="info"
    ) 