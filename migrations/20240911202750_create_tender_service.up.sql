CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
        CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'service_type') THEN
        CREATE TYPE service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
        CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT CHECK (LENGTH(description) <= 500),
    service_type service_type NOT NULL,
    status tender_status DEFAULT 'Created',
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1 CHECK (version >= 1),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    creator_username VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS tender_version (
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE, 
    name VARCHAR(100) NOT NULL,
    description TEXT CHECK (LENGTH(description) <= 500), 
    service_type service_type NOT NULL,
    status tender_status DEFAULT 'Created',
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1 CHECK (version >= 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    creator_username VARCHAR(50) NOT NULL,
    PRIMARY KEY (tender_id, version) 
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_status') THEN
        CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_author_type') THEN
        CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL CHECK (LENGTH(description) <= 500), 
    status bid_status NOT NULL DEFAULT 'Created',
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE, 
    author_type bid_author_type NOT NULL,
    author_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE, 
    version INT NOT NULL DEFAULT 1 CHECK (version >= 1), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tender_id, author_id)
);

CREATE TABLE IF NOT EXISTS bid_version (
    bid_id UUID NOT NULL REFERENCES bid(id) ON DELETE CASCADE, 
    name VARCHAR(100) NOT NULL,
    description TEXT CHECK (LENGTH(description) <= 500), 
    status bid_status NOT NULL, 
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE, 
    author_type bid_author_type NOT NULL, 
    author_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    version INT NOT NULL CHECK (version >= 1), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    PRIMARY KEY (bid_id, version) 
);

CREATE TABLE IF NOT EXISTS bid_decision (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL REFERENCES bid(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    decision VARCHAR(10) NOT NULL CHECK (decision IN ('Approved', 'Rejected')),  
    UNIQUE (bid_id, user_id) 
);


CREATE TABLE IF NOT EXISTS bid_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    bid_id UUID NOT NULL REFERENCES bid(id) ON DELETE CASCADE, 
    description TEXT CHECK (LENGTH(description) <= 500), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() 
);