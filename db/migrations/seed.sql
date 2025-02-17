SET
    statement_timeout = 0;

SET
    lock_timeout = 0;

SET
    idle_in_transaction_session_timeout = 0;

SET
    client_encoding = 'UTF8';

SET
    standard_conforming_strings = on;

SELECT
    pg_catalog.set_config ('search_path', '', false);

SET
    check_function_bodies = false;

SET
    xmloption = content;

SET
    client_min_messages = warning;

SET
    row_security = off;

SET
    default_tablespace = '';

SET
    default_table_access_method = heap;

CREATE TABLE public.attendance (
    id bigint NOT NULL,
    employee_id bigint NOT NULL,
    check_in timestamp
    with
        time zone NOT NULL,
        office_location_id bigint NOT NULL
);

ALTER TABLE public.attendance
ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.attendance_id_seq START
    WITH
        1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1
);

CREATE TABLE public.office_locations (
    id bigint NOT NULL,
    name text NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    employer_id bigint
);

ALTER TABLE public.office_locations
ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.office_locations_id_seq START
    WITH
        1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1
);

CREATE TABLE public.users (
    telegram_id integer NOT NULL,
    telegram_name text,
    type text,
    office_id integer
);

ALTER TABLE ONLY public.attendance ADD CONSTRAINT attendance_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.office_locations ADD CONSTRAINT office_locations_pkey PRIMARY KEY (id);

-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.users ADD CONSTRAINT users_pkey PRIMARY KEY (telegram_id);

--
-- Name: users users_telegram_id_key; Type: CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.users ADD CONSTRAINT users_telegram_id_key UNIQUE (telegram_id);

--
-- Name: attendance attendance_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.attendance ADD CONSTRAINT attendance_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.users (telegram_id) ON DELETE CASCADE;

--
-- Name: attendance attendance_office_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.attendance ADD CONSTRAINT attendance_office_location_id_fkey FOREIGN KEY (office_location_id) REFERENCES public.office_locations (id) ON DELETE CASCADE;

--
-- Name: office_locations office_locations_employer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.office_locations ADD CONSTRAINT office_locations_employer_id_fkey FOREIGN KEY (employer_id) REFERENCES public.users (telegram_id) ON DELETE SET NULL;

--
-- Name: users users_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--
ALTER TABLE ONLY public.users ADD CONSTRAINT users_office_id_fkey FOREIGN KEY (office_id) REFERENCES public.office_locations (id) ON DELETE SET NULL;
