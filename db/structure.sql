--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Debian 14.6-1.pgdg110+1)
-- Dumped by pg_dump version 14.6 (Homebrew)

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
    NEW.updated_at = now() at time zone 'utc';
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
-- Name: clients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.clients (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    email text DEFAULT ''::text NOT NULL,
    phone text DEFAULT ''::text NOT NULL,
    age integer DEFAULT 0 NOT NULL,
    weight integer DEFAULT 0 NOT NULL,
    height integer DEFAULT 0 NOT NULL,
    gender integer DEFAULT 0 NOT NULL,
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
    email text DEFAULT ''::text NOT NULL,
    phone text DEFAULT ''::text NOT NULL,
    age integer DEFAULT 0 NOT NULL,
    gender integer DEFAULT 0 NOT NULL,
    info text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL
);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


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
-- Name: auth_data_source_id_source_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX auth_data_source_id_source_type_idx ON public.auth_data USING btree (source_id, source_type);


--
-- Name: clients_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX clients_email_idx ON public.clients USING btree (email);


--
-- Name: nutritionists_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX nutritionists_email_idx ON public.nutritionists USING btree (email);


--
-- Name: auth_data update_auth_data_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_auth_data_updated_at BEFORE UPDATE ON public.auth_data FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: clients update_clients_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON public.clients FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: nutritionists update_nutritionists_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_nutritionists_updated_at BEFORE UPDATE ON public.nutritionists FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- PostgreSQL database dump complete
--

