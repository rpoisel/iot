SELECT
    d,
    COALESCE(excesspower, 0) AS excesspower,
    obtainedpower,
    COALESCE(solarpower - excesspower, 0) AS producedused,
    totalpower,
    solarpower,
    COALESCE((solarpower - excesspower) / solarpower, 1) AS ratio
FROM ( WITH deltas AS (
        SELECT
            modtime::date AS d,
            EXTRACT(EPOCH FROM (modtime - lag(modtime) OVER w)) AS timediff,
            solar,
            total,
            lag(solar) OVER w AS solarprev,
                lag(total) OVER w AS totalprev
                FROM
                    public.power
WINDOW w AS (ORDER BY modtime))
SELECT
    d,
    SUM(
        CASE WHEN solar > total THEN
        (solar + solarprev - (total + totalprev)) / 2 * timediff
        END) / (3600 * 1000) AS excesspower, SUM(
        CASE WHEN total > solar THEN
        (total + totalprev - (solar + solarprev)) / 2 * timediff
        END) / (3600 * 1000) AS obtainedpower,  SUM((total + totalprev) / 2 * timediff) / (3600 * 1000) AS totalpower, SUM((solar + solarprev) / 2 * timediff) / (3600 * 1000) AS solarpower
FROM
    deltas GROUP BY d ORDER BY d) AS powersums;

