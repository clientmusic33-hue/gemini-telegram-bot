#!/usr/bin/env python3
"""
Gemini Telegram Bot - Production-Ready AI Chatbot
Main application entry point
"""

import asyncio
import logging
import sys
from pathlib import Path

# Add src to path
sys.path.insert(0, str(Path(__file__).parent / "src"))

from config import load_config
from bot import GeminiTelegramBot
from utils.logging_config import setup_logging


async def main():
    """Main entry point"""
    # Load configuration
    config = load_config()
    
    # Setup logging
    setup_logging(config)
    logger = logging.getLogger(__name__)
    
    logger.info("=" * 60)
    logger.info(f"Starting {config.bot.name} v{config.bot.version}")
    logger.info("=" * 60)
    
    # Create and start bot
    bot = GeminiTelegramBot(config)
    
    try:
        await bot.start()
    except KeyboardInterrupt:
        logger.info("Received interrupt signal")
    except Exception as e:
        logger.error(f"Fatal error: {e}", exc_info=True)
        sys.exit(1)
    finally:
        await bot.stop()


if __name__ == "__main__":
    asyncio.run(main())
