-- live_data 테이블 생성
CREATE TABLE IF NOT EXISTS live_data (
                                         id SERIAL PRIMARY KEY,
                                         live_id INTEGER NOT NULL,
                                         live_title VARCHAR(255) NOT NULL,
    concurrent_user_count INTEGER NOT NULL DEFAULT 0,
    open_date TIMESTAMP WITH TIME ZONE NOT NULL,
                            adult BOOLEAN NOT NULL DEFAULT FALSE,
                            tags JSONB,
                            category_type VARCHAR(100) NOT NULL,
    live_category VARCHAR(100) NOT NULL,
    live_category_value VARCHAR(100) NOT NULL,
    channel_id VARCHAR(100) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    channel_image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
                            );

-- 인덱스 생성
CREATE INDEX IF NOT EXISTS idx_live_data_live_id ON live_data(live_id);
CREATE INDEX IF NOT EXISTS idx_live_data_channel_id ON live_data(channel_id);
CREATE INDEX IF NOT EXISTS idx_live_data_open_date ON live_data(open_date);
CREATE INDEX IF NOT EXISTS idx_live_data_created_at ON live_data(created_at);
