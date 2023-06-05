create DATABASE LTE1
    DATA DIRECTORY=’/mnt/disk1/db’
    INDEX DIRECTORY=’/mnt/disk2/index’;


create table info
(
    username  varchar(255) charset utf8 not null
        primary key,
    password  varchar(255) charset utf8 not null,
    level     int                       not null,
    confirmed tinyint(1)                null
);

create table tbC2I
(
    CITY        varchar(255) null,
    SCELL       varchar(255) not null,
    NCELL       varchar(255) not null,
    PrC2I9      float        null,
    C2I_Mean    float        null,
    Std         float        null,
    SampleCount float        null,
    WeightedC2I float        null,
    primary key (SCELL, NCELL)
);

create table tbC2I3
(
    a varchar(255) charset utf8 not null,
    b varchar(255) charset utf8 not null,
    c varchar(255) charset utf8 not null,
    primary key (a, b, c)
);

create table tbC2Inew
(
    SCELL    varchar(255) charset utf8 not null,
    NCELL    varchar(255) charset utf8 not null,
    RSRPmean float                     null,
    RSRPstd  float                     null,
    PrbC2I9  float                     null,
    PrbABS6  float                     null,
    primary key (SCELL, NCELL)
);

create index tbC2Inew_NCELL_SCELL_index
    on tbC2Inew (NCELL, SCELL);

create table tbCell
(
    CITY        varchar(255)   null,
    SECTOR_ID   varchar(255)   not null
        primary key,
    SECTOR_NAME varchar(255)   not null,
    ENODEBID    varchar(255)   not null,
    ENODEB_NAME varchar(255)   not null,
    EARFCN      int            not null,
    PCI         int            null,
    PSS         int            null,
    SSS         int            null,
    TAC         int            null,
    VENDOR      varchar(255)   null,
    LONGITUDE   decimal(10, 6) not null,
    LATITUDE    decimal(10, 6) not null,
    STYLE       varchar(255)   null,
    AZIMUTH     float          not null,
    HEIGHT      float          null,
    ELECTTILT   float          null,
    MECHTILT    float          null,
    TOTLETILT   float          null,
    constraint 检查_PCI
        check (`PCI` between 0 and 503)
);

create definer = root@localhost trigger insert_into_tbEnodeb_trigger
    after insert
    on tbCell
    for each row
BEGIN
    INSERT IGNORE INTO tbEnodeb (CITY, ENODEBID, ENODEB_NAME, VENDOR, LONGITUDE, LATITUDE, STYLE)
    VALUES (NEW.City, NEW.EnodebID, NEW.Enodeb_Name, NEW.Vendor, NEW.Longitude, NEW.Latitude, NEW.Style);
END;

create table tbEnodeb
(
    CITY        varchar(255) null,
    ENODEBID    int          not null
        primary key,
    ENODEB_NAME varchar(255) not null,
    VENDOR      varchar(255) null,
    LONGITUDE   float        not null,
    LATITUDE    float        not null,
    STYLE       varchar(255) null,
    constraint ENODEB_NAME
        unique (ENODEB_NAME)
);

create table tbKPI
(
    StartTime                    datetime      not null,
    ENODEB_NAME                  varchar(255)  null,
    SECTOR_DESCRIPTION           varchar(255)  not null,
    SECTOR_NAME                  varchar(255)  not null,
    RCCConnSUCC                  int           null,
    RCCConnATT                   int           null,
    RCCConnRATE                  decimal(7, 4) null,
    ERABConnSUCC                 int           null,
    ERABConnATT                  int           null,
    ERABConnRATE                 decimal(7, 4) null,
    ENODEB_ERABRel               int           null,
    SECTOR_ERABRel               int           null,
    ERABDropRateNew              decimal(7, 4) null,
    WirelessAccessRateAY         decimal(7, 4) null,
    ENODEB_UECtxRel              int           null,
    UEContextRel                 int           null,
    UEContextSUCC                int           null,
    WirelessDropRate             decimal(7, 4) null,
    ENODEB_InterFreqHOOutSUCC    int           null,
    ENODEB_InterFreqHOOutATT     int           null,
    ENODEB_IntraFreqHOOutSUCC    int           null,
    ENODEB_IntraFreqHOOutATT     int           null,
    ENODEB_InterFreqHOInSUCC     int           null,
    ENODEB_InterFreqHOInATT      int           null,
    ENODEB_IntraFreqHOInSUCC     int           null,
    ENODEB_IntraFreqHOInATT      int           null,
    ENODEB_HOInRate              decimal(7, 4) null,
    ENODEB_HOOutRate             decimal(7, 4) null,
    IntraFreqHOOutRateZSP        decimal(7, 4) null,
    InterFreqHOOutRateZSP        decimal(7, 4) null,
    HOSuccessRate                decimal(7, 4) null,
    PDCP_UplinkThroughput        bigint        null,
    PDCP_DownlinkThroughput      bigint        null,
    RRCRebuildReq                int           null,
    RRCRebuildRate               decimal(7, 4) null,
    SourceENB_IntraFreqHOOutSUCC int           null,
    SourceENB_InterFreqHOOutSUCC int           null,
    SourceENB_IntraFreqHOInSUCC  int           null,
    SourceENB_InterFreqHOInSUCC  int           null,
    ENODEB_HOOutSUCC             int           null,
    ENODEB_HOOutATT              int           null,
    primary key (SECTOR_NAME, StartTime),
    constraint FK_ENODEB_NAME
        foreign key (ENODEB_NAME) references tbEnodeb (ENODEB_NAME)
);

create table tbMROData
(
    TimeStamp         varchar(30) not null,
    ServingSector     varchar(50) not null,
    InterferingSector varchar(50) not null,
    LteScRSRP         float       null,
    LteNcRSRP         float       null,
    LteNcEarfcn       int         null,
    LteNcPci          smallint    null,
    primary key (TimeStamp, ServingSector, InterferingSector)
);

create table tbPRB
(
    StartTime          datetime     not null,
    ENODEB_NAME        varchar(255) null,
    SECTOR_DESCRIPTION varchar(255) not null,
    SECTOR_NAME        varchar(255) not null,
    PRB00              int          null,
    PRB01              int          null,
    PRB02              int          null,
    PRB03              int          null,
    PRB04              int          null,
    PRB05              int          null,
    PRB06              int          null,
    PRB07              int          null,
    PRB08              int          null,
    PRB09              int          null,
    PRB10              int          null,
    PRB11              int          null,
    PRB12              int          null,
    PRB13              int          null,
    PRB14              int          null,
    PRB15              int          null,
    PRB16              int          null,
    PRB17              int          null,
    PRB18              int          null,
    PRB19              int          null,
    PRB20              int          null,
    PRB21              int          null,
    PRB22              int          null,
    PRB23              int          null,
    PRB24              int          null,
    PRB25              int          null,
    PRB26              int          null,
    PRB27              int          null,
    PRB28              int          null,
    PRB29              int          null,
    PRB30              int          null,
    PRB31              int          null,
    PRB32              int          null,
    PRB33              int          null,
    PRB34              int          null,
    PRB35              int          null,
    PRB36              int          null,
    PRB37              int          null,
    PRB38              int          null,
    PRB39              int          null,
    PRB40              int          null,
    PRB41              int          null,
    PRB42              int          null,
    PRB43              int          null,
    PRB44              int          null,
    PRB45              int          null,
    PRB46              int          null,
    PRB47              int          null,
    PRB48              int          null,
    PRB49              int          null,
    PRB50              int          null,
    PRB51              int          null,
    PRB52              int          null,
    PRB53              int          null,
    PRB54              int          null,
    PRB55              int          null,
    PRB56              int          null,
    PRB57              int          null,
    PRB58              int          null,
    PRB59              int          null,
    PRB60              int          null,
    PRB61              int          null,
    PRB62              int          null,
    PRB63              int          null,
    PRB64              int          null,
    PRB65              int          null,
    PRB66              int          null,
    PRB67              int          null,
    PRB68              int          null,
    PRB69              int          null,
    PRB70              int          null,
    PRB71              int          null,
    PRB72              int          null,
    PRB73              int          null,
    PRB74              int          null,
    PRB75              int          null,
    PRB76              int          null,
    PRB77              int          null,
    PRB78              int          null,
    PRB79              int          null,
    PRB80              int          null,
    PRB81              int          null,
    PRB82              int          null,
    PRB83              int          null,
    PRB84              int          null,
    PRB85              int          null,
    PRB86              int          null,
    PRB87              int          null,
    PRB88              int          null,
    PRB89              int          null,
    PRB90              int          null,
    PRB91              int          null,
    PRB92              int          null,
    PRB93              int          null,
    PRB94              int          null,
    PRB95              int          null,
    PRB96              int          null,
    PRB97              int          null,
    PRB98              int          null,
    PRB99              int          null,
    primary key (SECTOR_NAME, StartTime)
);

insert into info (username, password, level, confirmed) values ('admin','$2a$10$cLT3CBOwUzp4avQYd0uNeegm6b.tXktbRwYz.VIduCJOXVWQbhVOu',0,true);
