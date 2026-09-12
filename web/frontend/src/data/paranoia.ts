/**
 * 《妄想症 Paranoia》系列资料
 *
 * 致敬对象：DELA 作曲 / 雨狸作词 / Lune 等人曲绘的 VOCALOID 中文原创系列曲
 * （2015.11.28《一重加害》— 2017.03.17《九重现实》，共 475 天；
 * 2025.11.28 发布十周年纪念曲《十重告别》）。
 *
 * 曲绘、歌词与设定的著作权均归原作者所有，此处仅作主题化的介绍与引用。
 */

import ch01 from '@/assets/paranoia/ch01.jpg'
import ch02 from '@/assets/paranoia/ch02.jpg'
import ch03 from '@/assets/paranoia/ch03.jpg'
import ch04 from '@/assets/paranoia/ch04.jpg'
import ch05 from '@/assets/paranoia/ch05.jpg'
import ch06 from '@/assets/paranoia/ch06.jpg'
import ch07 from '@/assets/paranoia/ch07.jpg'
import ch08 from '@/assets/paranoia/ch08.jpg'
import ch09 from '@/assets/paranoia/ch09.jpg'
import ch10 from '@/assets/paranoia/ch10.jpg'

/** 章：乐曲在系列中的位置 */
export type ParanoiaArc = '序曲' | '坠落之章' | '上升之章' | '十周年'

export interface ParanoiaChapter {
  /** 重数，0 表示序曲 */
  no: number
  /** 中文数字，如「一」 */
  numeral: string
  /** 曲名 */
  title: string
  arc: ParanoiaArc
  /** 乐谱记号：坠落之章为 b（降号），上升之章为 #（升号） */
  sigil: string
  /** 投稿日期 */
  date: string
  /** 演唱 */
  singers: string
  /** 视角角色 */
  voices: string
  /** 一句话剧情 */
  summary: string
  /** 章节简介 */
  detail: string
  /** 歌词摘录 */
  lyric: string
  /** 对应的妄想症条目 */
  delusion: string
  /** 临床描述 */
  note: string
  image: string
}

/** 三重加害者、守护者、背叛者与第四面墙外之人 */
export interface ParanoiaCharacter {
  name: string
  /** 虚拟歌手 */
  vocal: string
  role: string
  /** 主题色 */
  color: string
  /** 第二主题色，用于言语这种双重身份的角色 */
  colorAlt?: string
  /** 一句属于她的话 */
  line: string
}

/** 妄想档案：系列中逐章出现的临床条目 */
export interface ParanoiaDelusion {
  index: string
  term: string
  alias: string
  note: string
  chapter: string
}

export const paranoiaChapters: ParanoiaChapter[] = [
  {
    no: 1,
    numeral: '一',
    title: '一重加害',
    arc: '坠落之章',
    sigil: 'b',
    date: '2015.11.28',
    singers: '洛天依 / 言和',
    voices: '被害者泠珞 · 加害者颜语',
    summary: '熟悉的街道上多了一个影子，被害妄想从此降临。',
    detail:
      '系列的开篇。以被害妄想为主题，由被害者泠珞、加害者颜语与背叛者零羽三人展开，剧情扑朔迷离，在留下巨大悬念的同时，也为最终的结局埋下伏笔。',
    lyric: '别对我说 你身后是谁在揣度\n别靠近我 你这无锁链的恶魔',
    delusion: '被害妄想',
    note: '患者往往处于恐惧状态而胡乱推理和分析，坚持自己受到迫害或伤害，如言语上的针对、嘲弄、跟踪及监听等。',
    image: ch01,
  },
  {
    no: 2,
    numeral: '二',
    title: '二重变革',
    arc: '坠落之章',
    sigil: 'b',
    date: '2015.12.12',
    singers: '洛天依 / 言和',
    voices: '被害者泠珞 · 加害者颜语',
    summary: '凶手遵从风声鹤唳者的妄想而来，自虚无中诞生。',
    detail:
      '第一作的故事后续与不同视角描述。被害者与加害者之间的斗争正式开始：「我正是接受你『会被杀死』的强烈夙愿而自虚无诞生。」',
    lyric: '这一重任扭曲泛滥的世界\n它可悲的明天就由我来终结',
    delusion: '影响妄想 · 逆',
    note: '患者认为自己的思维、情感与意志行为可以影响并支配外界的现实，即「心想事成」。',
    image: ch02,
  },
  {
    no: 3,
    numeral: '三',
    title: '三重爱恋',
    arc: '坠落之章',
    sigil: 'b',
    date: '2016.01.22',
    singers: '洛天依 / 言和',
    voices: '被害者泠珞 · 守护者颜语 · 潜意识泠珞',
    summary: '杀死幻想怪物之后，她被裁定防卫过当，于是开始了一场妄想中的恋爱。',
    detail:
      '盲目生活像陀螺一般旋转不停，动荡的视线引发不安。在梦魇深堕之处，她开始相信自己终于遇到了那个人。',
    lyric: '约会在吊死鬼开的乐园\n看糖果和孤独成就了饕餮',
    delusion: '情爱妄想',
    note: '又称 de Clerambault 症候群。病人会以为自己正和某人恋爱，或另一个人深爱着自己；幻想中的恋人甚至只是一个「魅影」。',
    image: ch03,
  },
  {
    no: 4,
    numeral: '四',
    title: '四重罪孽',
    arc: '坠落之章',
    sigil: 'b',
    date: '2016.03.26',
    singers: '洛天依 / 言和 / 乐正绫',
    voices: '潜意识泠珞 · 守护者颜语 · 背叛者零羽',
    summary: '潜意识泠珞觉醒，察觉一切皆妄想，并宣告自己犯下的重重罪孽。',
    detail:
      '表现挣扎的一章。加害者的复活、守护者准备为爱牺牲的觉悟、背叛者扑朔迷离的死因，都替之后的剧情留下伏笔。',
    lyric: '我犯下了重重罪孽 我到底是什么存在\n我是最恐怖的危害 能否从此逃开',
    delusion: '自罪妄想',
    note: '又称罪恶妄想。患者毫无根据地认为自己犯了严重错误和罪行，或将一件小事的错误无限夸大，以至于拒食或要求赎罪。',
    image: ch04,
  },
  {
    no: 5,
    numeral: '五',
    title: '五重空洞',
    arc: '坠落之章',
    sigil: 'b',
    date: '2016.04.28',
    singers: '洛天依 / 乐正绫',
    voices: '潜意识泠珞 · 背叛者零羽',
    summary: '妄想世界的夹缝里，彼岸花田中埋葬着背叛者零羽。',
    detail:
      '潜意识泠珞在彼岸花田中见到了埋葬在记忆深处的零羽，过去的种种记忆漫上心头。故事的真相揭起一角，过去与未来依旧扑朔迷离。',
    lyric: '你停留在我无法追逐的记忆彼岸\n永远不会再回来',
    delusion: '内心被揭露感',
    note: '又称被洞悉感。患者认为其内心的想法未经表达，别人就知道了；严重时表层意识可能被穿透，导致潜意识暴露在外。',
    image: ch05,
  },
  {
    no: 6,
    numeral: '六',
    title: '六重不忠',
    arc: '上升之章',
    sigil: '#',
    date: '2016.08.05',
    singers: '洛天依 / 言和 / 乐正绫',
    voices: '加害者 · 守护者 · 潜意识泠珞',
    summary: '守护者先被策反，又因誓言死在加害者刀下，系列由此转入上升之章。',
    detail:
      '从本曲开始，系列进入 # 之章（上升之章）。她先期待、她先许愿，身为镜像怎能妥协；而誓言与爱，最终以死亡兑现。',
    lyric: '再见吧终会相见 莫为此放弃明天\n在绝望之前 将我的誓言默念',
    delusion: '嫉妒妄想综合征',
    note: '又称奥赛罗综合征。患者坚信所爱之人不忠，即使有确切证据证明其并未背叛，怀疑的依据也捕风捉影。',
    image: ch06,
  },
  {
    no: 7,
    numeral: '七',
    title: '七重痼病',
    arc: '上升之章',
    sigil: '#',
    date: '2016.12.16',
    singers: '洛天依 / 言和 / 乐正绫',
    voices: '被害者泠珞 · 加害者颜语 · 零羽',
    summary: '加害者与潜意识同归于尽，她在白日的中心失去焦距。',
    detail:
      '第六作的后续。红石蒜在无人之地盛开，妄议着未来于是溃败。「如果谎言重复了千遍依然无法成为现实，那谎言应该如何生存下去呢？」',
    lyric: '世界已发出信号 停止临终前无止境祈祷\n我接受她不在的未来',
    delusion: '疑病妄想',
    note: '患者毫无根据地坚信自己患了某种严重躯体疾病或不治之症。严重时称为虚无幻想——「本人已经不存在，只剩下一个躯壳空壳了」。',
    image: ch07,
  },
  {
    no: 8,
    numeral: '八',
    title: '八重回归',
    arc: '上升之章',
    sigil: '#',
    date: '2017.03.10',
    singers: '洛天依 / 乐正绫',
    voices: '泠珞 · 零羽',
    summary: '一切妄想破灭之后回归现实，同时与过往和梦想告别。',
    detail:
      '睁开眼重新堕入平凡的世界，辗转枯燥的时间一天又一天。所有羁绊终归于平淡，誓言悉数臣服现实的安排——「无意义乃这个世界存在的唯一意义。」',
    lyric: '再见 再不见\n无比真实过的一切',
    delusion: '回归空白',
    note: '现实本身就需要和妄想区分开来。越痛苦，翅膀越充血，它们就越美丽——然后死亡就突然降临，只是轻飘飘地、轻飘飘地。',
    image: ch08,
  },
  {
    no: 9,
    numeral: '九',
    title: '九重现实',
    arc: '上升之章',
    sigil: '#',
    date: '2017.03.17',
    singers: '洛天依 / 乐正绫 · 言和和声',
    voices: '泠珞 · 零羽',
    summary: '完结作。故事的真正脉络，以及带着希望的未来。',
    detail:
      '从「与其生离，不如死别」到「面对凶险的今后」。不是自作多情，也不是阳光明媚，仅仅是空白的现实而已。白色，也许是可以着色的，也许是被擦除了色彩之后的。',
    lyric: '你的声音击破梦魇一重一重\n自最遥远的彼方呼唤我',
    delusion: '现实',
    note: '「我所描述的一切或许并非真实，然而我的感受绝无虚假。」绝无虚假。',
    image: ch09,
  },
  {
    no: 10,
    numeral: '十',
    title: '十重告别',
    arc: '十周年',
    sigil: '✳',
    date: '2025.11.28',
    singers: '洛天依 / 言和 / 乐正绫 / 双狐座',
    voices: '第四面墙外之人',
    summary: '第四面墙外之人回望整个系列，把苦痛与欢愉收纳成歌颂明天的力量。',
    detail:
      '十周年纪念曲。打破第四面墙的元叙事视角，带领角色、创作者与听众共同回望这一切：告别或是全新的起点，全在他者的一念之间。',
    lyric: '无论再见或再也不见 我们心声已然相遇见\n比十年更早之前 我便歌颂着明天',
    delusion: '告别',
    note: '载体或许会消失，但我们相遇过的回忆不会。无论是十年前还是十年后，都希望你能拥有属于自己的幸福。',
    image: ch10,
  },
]

export const paranoiaCharacters: ParanoiaCharacter[] = [
  {
    name: '泠珞',
    vocal: '洛天依',
    role: '被害者',
    color: '#66ccff',
    line: '不承认那是命运，命运是被神抛弃的棋局。',
  },
  {
    name: '颜语',
    vocal: '言和',
    role: '加害者 / 守护者',
    color: '#f2d16b',
    colorAlt: '#7ee8e0',
    line: '如果就此离开，还要怎么加害。',
  },
  {
    name: '零羽',
    vocal: '乐正绫',
    role: '背叛者',
    color: '#e0245e',
    line: '背对着我一跃而下的女孩，在呼唤。',
  },
  {
    name: '双狐座',
    vocal: '第四面墙外之人',
    role: '元叙事视角',
    color: '#b98cff',
    line: '告别或是全新的起点，全在他者的一念之间。',
  },
]

export const paranoiaDelusions: ParanoiaDelusion[] = [
  {
    index: '01',
    term: '被害妄想',
    alias: 'delusion of persecution',
    note: '处于恐惧状态而胡乱推理和分析，坚持自己受到迫害或伤害，如言语上的针对、嘲弄、跟踪及监听。',
    chapter: '一重加害',
  },
  {
    index: '02',
    term: '影响妄想',
    alias: 'delusion of influence',
    note: '认为自己的思维、情感、意志行为活动受到外界某种力量的支配、控制、操纵，不能自主。',
    chapter: '二重变革',
  },
  {
    index: '03',
    term: '影响妄想 · 逆',
    alias: 'reverse influence',
    note: '认为自己的思维与意志不论多么荒谬，都可以影响并支配外界的现实，即「心想事成」。',
    chapter: '二重变革',
  },
  {
    index: '04',
    term: '情爱妄想',
    alias: "de Clerambault's syndrome",
    note: '以为自己正和某人恋爱，或另一个人深爱着自己；幻想中的恋人甚至只是一个现实中并不存在的「魅影」。',
    chapter: '三重爱恋',
  },
  {
    index: '05',
    term: '自罪妄想',
    alias: 'delusion of guilt',
    note: '又称罪恶妄想。毫无根据地认为自己犯了严重错误和罪行，罪大恶极、死有余辜，应该受到惩罚。',
    chapter: '四重罪孽',
  },
  {
    index: '06',
    term: '内心被揭露感',
    alias: 'experience of being revealed',
    note: '又称被洞悉感。认为内心的想法未经表达别人就知道了；严重时表层意识被穿透，潜意识暴露在外。',
    chapter: '五重空洞',
  },
  {
    index: '07',
    term: '投射性认同',
    alias: 'projective identification',
    note: '诱导他人以一种限定的方式行动或做出反应的人际行为模式，接受者往往被迫成为投射者感受的储存室。',
    chapter: '五重空洞',
  },
  {
    index: '08',
    term: '嫉妒妄想',
    alias: 'Othello syndrome',
    note: '坚信所爱之人不贞，即使有确切证据也丝毫不动摇；思维推理荒谬，怀疑的依据捕风捉影。',
    chapter: '六重不忠',
  },
  {
    index: '09',
    term: '疑病妄想',
    alias: 'hypochondriacal delusion',
    note: '毫无根据地坚信自己患了某种严重躯体疾病或不治之症，反复检查也不能纠正其歪曲的信念。',
    chapter: '七重痼病',
  },
  {
    index: '10',
    term: '虚无幻想',
    alias: 'nihilistic delusion',
    note: '疑病妄想的重度形态：认为「本人已经不存在，只剩下一个躯壳空壳了」。',
    chapter: '七重痼病',
  },
]

/** 《九重现实》中对妄想的五条描述 */
export const paranoiaAxioms: string[] = [
  '妄想是一个根本性障碍。',
  '是慢性和终身频发的。',
  '妄想有逻辑构造，并且内部成立。',
  '妄想和普通逻辑推理不冲突，而且没有一般的行为异常。',
  '自身经历了对自我参考的高度理解——对别人来说不重要的事件，对他们来说极为重要。',
]

/** 主题化引文，用于顶栏滚动字幕与仪表盘签文 */
export const paranoiaQuotes: { text: string; from: string }[] = [
  { text: '谎言重复一千次就变成真理。所以，现实也是这样诞生的。', from: '《一重加害》' },
  { text: '梦魇中是谁背对着我纵身一跃，没说再见，也没有遗言。', from: '《一重加害》' },
  { text: '这一重任扭曲泛滥的世界，它可悲的明天就由我来终结。', from: '《二重变革》' },
  { text: '约会在吊死鬼开的乐园，看糖果和孤独成就了饕餮。', from: '《三重爱恋》' },
  { text: '我犯下了重重罪孽，我到底是什么存在。', from: '《四重罪孽》' },
  { text: '你停留在我无法追逐的记忆彼岸，永远不会再回来。', from: '《五重空洞》' },
  { text: '再见吧终会相见，莫为此放弃明天。', from: '《六重不忠》' },
  { text: '如果谎言重复了千遍依然无法成为现实，那谎言应该如何生存下去呢？', from: '《七重痼病》' },
  { text: '无意义乃这个世界存在的唯一意义。', from: '《八重回归》' },
  { text: '我所描述的一切或许并非真实，然而我的感受绝无虚假。', from: '《九重现实》' },
  { text: '比十年更早之前，我便歌颂着明天。', from: '《十重告别》' },
]

/** 曲目列表：同系列其他作品 */
export const paranoiaSeries: { group: string; items: string[] }[] = [
  { group: '序曲', items: ['《NOVA》'] },
  {
    group: '正篇',
    items: [
      '《一重加害》',
      '《二重变革》',
      '《三重爱恋》',
      '《四重罪孽》',
      '《五重空洞》',
      '《六重不忠》',
      '《七重痼病》',
      '《八重回归》',
      '《九重现实》',
    ],
  },
  {
    group: '角色歌',
    items: ['《三重恋爱》', '《泠重乞愿》', '《零重祈愿》'],
  },
  {
    group: '外篇',
    items: ['《自攻自受》', '《妄想 Paranoia》', '《妄想 Reality》', '《十重告别》'],
  },
]

/**
 * 侧边栏章节谱：把应用的每个功能对应到系列的一重上。
 * 坠落之章用 b（降号），上升之章用 #（升号）——这也是原作分章的方式。
 */
export interface ParanoiaChapterMark {
  /** 章节记号 */
  sigil: string
  /** 展示用的重数 */
  label: string
  /** 曲名 */
  title: string
}

export const paranoiaNavMarks: Record<string, ParanoiaChapterMark> = {
  dashboard: { sigil: '✳', label: '序曲', title: 'NOVA' },
  timers: { sigil: 'b', label: '一重', title: '一重加害' },
  rules: { sigil: 'b', label: '二重', title: '二重变革' },
  accounts: { sigil: 'b', label: '三重', title: '三重爱恋' },
  gateways: { sigil: 'b', label: '四重', title: '四重罪孽' },
  'gateway-rules': { sigil: 'b', label: '五重', title: '五重空洞' },
  logs: { sigil: '#', label: '六重', title: '六重不忠' },
  settings: { sigil: '#', label: '七重', title: '七重痼病' },
  logout: { sigil: '#', label: '八重', title: '八重回归' },
  paranoia: { sigil: '#', label: '九重', title: '九重现实' },
}

export const paranoiaCredits = {
  title: '妄想症 Paranoia',
  span: '2015.11.28 — 2017.03.17 · 共 475 天',
  staff: [
    { role: '作编曲', name: 'DELA' },
    { role: '作词 / 策划', name: '雨狸' },
    { role: '曲绘', name: 'Lune、YAL、雨湘雪' },
    { role: 'PV', name: '一折起售、owen-z' },
    { role: '混音', name: 'POiSON' },
    { role: '调校', name: '花儿不哭、Digger、河谐、OQQ、流绪' },
  ],
  notice:
    '本主题为界面外观层面的同人致敬，非官方合作。「妄想症 Paranoia」系列曲目、曲绘、歌词与设定的著作权均归原作者所有；此处仅作介绍与引用，不作任何商业用途。',
}

export function chapterByNo(no: number): ParanoiaChapter | undefined {
  return paranoiaChapters.find((item) => item.no === no)
}

/** 按日期轮换一条引文，让「今日签文」每天都不一样 */
export function todaysQuote(date = new Date()) {
  const start = new Date(date.getFullYear(), 0, 0)
  const dayOfYear = Math.floor((date.getTime() - start.getTime()) / 86400000)
  return paranoiaQuotes[dayOfYear % paranoiaQuotes.length]
}
