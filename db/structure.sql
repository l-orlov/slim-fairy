--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Debian 14.6-1.pgdg110+1)
-- Dumped by pg_dump version 14.7 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: set_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now() AT TIME ZONE 'utc';
    RETURN NEW;
END;
$$;


SET default_table_access_method = heap;

--
-- Name: auth_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.auth_data (
    source_id uuid NOT NULL,
    source_type text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL
);


--
-- Name: chat_bot_dialogs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_bot_dialogs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_telegram_id bigint NOT NULL,
    kind text NOT NULL,
    status integer NOT NULL,
    data jsonb,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL
);


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: nutritionists; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.nutritionists (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    email text,
    phone text,
    telegram_id bigint,
    age integer,
    gender integer,
    info text,
    created_by text,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    email text,
    phone text,
    telegram_id bigint,
    age integer,
    weight integer,
    height integer,
    gender integer,
    physical_activity_level integer,
    created_by text,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL
);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: chat_bot_dialogs chat_bot_dialogs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_bot_dialogs
    ADD CONSTRAINT chat_bot_dialogs_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: nutritionists nutritionists_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nutritionists
    ADD CONSTRAINT nutritionists_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: auth_data_source_id_source_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX auth_data_source_id_source_type_idx ON public.auth_data USING btree (source_id, source_type);


--
-- Name: chat_bot_dialogs_user_telegram_id_kind_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX chat_bot_dialogs_user_telegram_id_kind_status_idx ON public.chat_bot_dialogs USING btree (user_telegram_id, kind, status);


--
-- Name: nutritionists_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX nutritionists_email_idx ON public.nutritionists USING btree (email);


--
-- Name: nutritionists_phone_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX nutritionists_phone_idx ON public.nutritionists USING btree (phone);


--
-- Name: nutritionists_telegram_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX nutritionists_telegram_id_idx ON public.nutritionists USING btree (telegram_id);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: users_phone_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_phone_idx ON public.users USING btree (phone);


--
-- Name: users_telegram_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_telegram_id_idx ON public.users USING btree (telegram_id);


--
-- Name: chat_bot_dialogs chat_bot_dialogs_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER chat_bot_dialogs_updated_at BEFORE UPDATE ON public.chat_bot_dialogs FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: auth_data update_auth_data_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_auth_data_updated_at BEFORE UPDATE ON public.auth_data FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: nutritionists update_nutritionists_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_nutritionists_updated_at BEFORE UPDATE ON public.nutritionists FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- PostgreSQL database dump complete
--

