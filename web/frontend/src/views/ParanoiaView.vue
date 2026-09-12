<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  paranoiaAxioms,
  paranoiaChapters,
  paranoiaCharacters,
  paranoiaCredits,
  paranoiaDelusions,
  paranoiaSeries,
  type ParanoiaChapter,
} from '@/data/paranoia'
import { useIsMobile } from '@/composables/useBreakpoint'

/**
 * 妄想症 · 九重档案
 *
 * 把 ./art/paranoia 下的曲绘与设定文本搬进前端：
 * 九重曲目、三重加害者与守护者、十份妄想病历，以及整个系列的落款。
 */

const isMobile = useIsMobile()

const filters = [
  { id: 'all', label: '全部九重' },
  { id: 'b', label: 'b 坠落之章' },
  { id: '#', label: '# 上升之章' },
] as const

const filter = ref<(typeof filters)[number]['id']>('all')

const visibleChapters = computed(() =>
  filter.value === 'all'
    ? paranoiaChapters
    : paranoiaChapters.filter((chapter) => chapter.sigil === filter.value),
)

const active = ref<ParanoiaChapter | null>(null)

const dialogOpen = computed({
  get: () => active.value !== null,
  set: (value: boolean) => {
    if (!value) active.value = null
  },
})

/** 首屏叠放的三张曲绘 */
const heroArt = [paranoiaChapters[4], paranoiaChapters[8], paranoiaChapters[2]]
</script>

<template>
  <div class="pa-page">
    <!-- ================= 卷首 ================= -->
    <section class="pa-hero">
      <div class="pa-hero__text gb-rise">
        <div class="pa-hero__eyebrow">
          <span class="pa-mark"><em>✳</em> 致敬篇</span>
          <span class="muted">妄想症 Paranoia · 九重档案</span>
        </div>

        <h1 class="pa-title">妄想症</h1>
        <p class="pa-title-sub">Paranoia</p>

        <p class="pa-hero__lede">
          被害者泠珞、加害者与守护者颜语、背叛者零羽——<br />
          九重正篇，四百七十五天，一场层层叠叠的戏剧。
        </p>

        <div class="pa-quote pa-hero__quote">
          谎言重复一千次就变成真理。<br />
          所以，现实也是这样诞生的。
          <span class="pa-hero__quote-from">——《一重加害》</span>
        </div>

        <div class="pa-hero__note">
          {{ paranoiaCredits.notice }}
        </div>
      </div>

      <div class="pa-hero__art gb-rise" aria-hidden="true">
        <div v-for="(item, i) in heroArt" :key="item.title" class="pa-hero__card" :class="`is-${i}`">
          <div class="pa-film"><img :src="item.image" alt="" /></div>
          <span class="pa-hero__card-name">{{ item.numeral }}重 · {{ item.title }}</span>
        </div>
      </div>
    </section>

    <!-- ================= 时间轴 ================= -->
    <div class="pa-timeline gb-rise">
      <div class="pa-timeline__node">
        <span class="pa-timeline__date mono">2015.11.28</span>
        <span class="pa-timeline__label">《一重加害》投稿</span>
      </div>
      <span class="pa-timeline__line"><i /></span>
      <div class="pa-timeline__node is-center">
        <span class="pa-timeline__date mono">475 DAYS</span>
        <span class="pa-timeline__label">坠落之章 b → 上升之章 #</span>
      </div>
      <span class="pa-timeline__line"><i /></span>
      <div class="pa-timeline__node is-end">
        <span class="pa-timeline__date mono">2025.11.28</span>
        <span class="pa-timeline__label">十周年《十重告别》</span>
      </div>
    </div>

    <!-- ================= 角色 ================= -->
    <section class="pa-section">
      <header class="pa-section__head">
        <h2 class="pa-section__title">登场人物</h2>
        <span class="pa-section__sub mono">CAST / 三重妄想</span>
      </header>

      <div class="pa-cast">
        <article
          v-for="cast in paranoiaCharacters"
          :key="cast.name"
          class="pa-cast__card gb-rise"
          :style="{ '--pa-cast': cast.color, '--pa-cast-alt': cast.colorAlt ?? cast.color }"
        >
          <div class="pa-cast__head">
            <span class="pa-cast__avatar">{{ cast.name[0] }}</span>
            <div>
              <div class="pa-cast__name">{{ cast.name }}</div>
              <div class="pa-cast__vocal mono">{{ cast.vocal }}</div>
            </div>
          </div>
          <div class="pa-cast__role">{{ cast.role }}</div>
          <p class="pa-cast__line">{{ cast.line }}</p>
        </article>
      </div>
    </section>

    <!-- ================= 九重曲目 ================= -->
    <section class="pa-section">
      <header class="pa-section__head">
        <h2 class="pa-section__title">九重曲目</h2>
        <span class="pa-section__sub mono">b FLAT / # SHARP</span>
      </header>
      <p class="pa-section__desc">
        前半用降号 b 记作「坠落之章」，后半用升号 # 记作「上升之章」——同一个故事的两半。
        点击任意一张卷宗，可以看到该重的剧情、歌词与对应的妄想条目。
      </p>

      <div class="pa-filters">
        <button
          v-for="item in filters"
          :key="item.id"
          type="button"
          class="pa-filter"
          :class="{ 'is-active': filter === item.id }"
          @click="filter = item.id"
        >
          {{ item.label }}
        </button>
      </div>

      <div class="pa-chapters">
        <button
          v-for="chapter in visibleChapters"
          :key="chapter.no"
          type="button"
          class="pa-chapter gb-rise"
          @click="active = chapter"
        >
          <div class="pa-film pa-chapter__art">
            <img :src="chapter.image" :alt="`《${chapter.title}》曲绘`" loading="lazy" />
            <span class="pa-chapter__sigil mono">{{ chapter.sigil }}{{ chapter.numeral }}</span>
          </div>
          <div class="pa-chapter__body">
            <div class="pa-chapter__meta">
              <span>{{ chapter.arc }}</span>
              <span class="mono">{{ chapter.date }}</span>
            </div>
            <h3 class="pa-chapter__title">《{{ chapter.title }}》</h3>
            <p class="pa-chapter__summary">{{ chapter.summary }}</p>
            <p class="pa-chapter__lyric">{{ chapter.lyric.split('\n')[0] }}</p>
            <span class="pa-chapter__more">展开卷宗 →</span>
          </div>
        </button>
      </div>
    </section>

    <!-- ================= 妄想档案 ================= -->
    <section class="pa-section">
      <header class="pa-section__head">
        <h2 class="pa-section__title">妄想档案</h2>
        <span class="pa-section__sub mono">CASE FILE / 10</span>
      </header>
      <p class="pa-section__desc">
        系列每一重都会在 PV 里夹一段临床描述。这些条目既是故事的注脚，也是角色们真正的处境。
      </p>

      <div class="pa-records">
        <article v-for="item in paranoiaDelusions" :key="item.index" class="pa-record gb-rise">
          <div class="pa-record__top">
            <span class="pa-record__index mono">{{ item.index }}</span>
            <span class="pa-record__chapter">{{ item.chapter }}</span>
          </div>
          <h3 class="pa-record__term">{{ item.term }}</h3>
          <p class="pa-record__alias mono">{{ item.alias }}</p>
          <p class="pa-record__note">{{ item.note }}</p>
        </article>
      </div>
    </section>

    <!-- ================= 五条描述 ================= -->
    <section class="pa-section pa-section--split">
      <div class="pa-split__col">
        <header class="pa-section__head">
          <h2 class="pa-section__title">关于妄想</h2>
          <span class="pa-section__sub mono">《九重现实》</span>
        </header>
        <ol class="pa-axioms">
          <li v-for="(axiom, i) in paranoiaAxioms" :key="axiom">
            <span class="mono">{{ String(i + 1).padStart(2, '0') }}</span>
            <p>{{ axiom }}</p>
          </li>
        </ol>
      </div>

      <div class="pa-split__col">
        <header class="pa-section__head">
          <h2 class="pa-section__title">系列曲目</h2>
          <span class="pa-section__sub mono">DISCOGRAPHY</span>
        </header>
        <div v-for="group in paranoiaSeries" :key="group.group" class="pa-disc">
          <div class="pa-disc__group">{{ group.group }}</div>
          <div class="pa-disc__items">
            <span v-for="song in group.items" :key="song" class="pa-disc__item">{{ song }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ================= 落款 ================= -->
    <footer class="pa-colophon">
      <div class="pa-rule"><span class="mono">谢幕 · 撒花</span></div>

      <p class="pa-colophon__quote">
        载体或许会消失，但我们相遇过的回忆不会。<br />
        无论是十年前还是十年后，都希望你能拥有属于自己的幸福。
      </p>

      <div class="pa-colophon__staff">
        <span v-for="person in paranoiaCredits.staff" :key="person.role" class="pa-colophon__item">
          <em>{{ person.role }}</em>{{ person.name }}
        </span>
      </div>

      <p class="pa-colophon__span">{{ paranoiaCredits.span }}</p>
      <p class="pa-colophon__notice">{{ paranoiaCredits.notice }}</p>
    </footer>

    <!-- ================= 卷宗详情 ================= -->
    <el-dialog
      v-model="dialogOpen"
      :title="active ? `《${active.title}》` : ''"
      :width="isMobile ? '92vw' : '720px'"
      align-center
      class="pa-dialog"
    >
      <div v-if="active" class="pa-dossier">
        <div class="pa-film pa-dossier__art">
          <img :src="active.image" :alt="`《${active.title}》曲绘`" />
          <span class="pa-dossier__sigil mono">{{ active.sigil }}{{ active.numeral }}</span>
        </div>

        <div class="pa-dossier__meta">
          <span class="pa-mark"><em>{{ active.sigil }}</em>{{ active.arc }} · {{ active.numeral }}重</span>
          <span class="muted mono">{{ active.date }}</span>
          <span class="muted">演唱：{{ active.singers }}</span>
        </div>

        <p class="pa-dossier__voices">视角：{{ active.voices }}</p>
        <p class="pa-dossier__detail">{{ active.detail }}</p>
        <p class="pa-quote pa-dossier__lyric">{{ active.lyric }}</p>

        <div class="pa-dossier__case">
          <div class="pa-dossier__case-term">
            <span class="mono">妄想条目</span>
            <strong>{{ active.delusion }}</strong>
          </div>
          <p>{{ active.note }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.pa-page {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 34px;
}

/* ---------- 卷首 ---------- */

.pa-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 0.92fr);
  align-items: center;
  gap: 32px;
  padding: 30px 30px 34px;
  border: 1px solid var(--pa-line);
  border-radius: 8px;
  background:
    radial-gradient(120% 100% at 0% 0%, rgb(224 36 94 / 0.12), transparent 58%),
    linear-gradient(180deg, var(--gb-card), transparent);
  overflow: hidden;
}

.pa-hero__eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.pa-title {
  margin: 0;
  font-family: var(--pa-serif);
  font-size: clamp(40px, 6vw, 60px);
  font-weight: 800;
  line-height: 1;
  letter-spacing: 0.06em;
  color: var(--el-text-color-primary);
}

.pa-title-sub {
  margin: 6px 0 18px;
  font-family: var(--pa-mono);
  font-size: 13px;
  letter-spacing: 0.42em;
  text-transform: uppercase;
  color: var(--gb-primary);
}

.pa-hero__lede {
  margin: 0 0 20px;
  font-size: 14px;
  line-height: 2;
  color: var(--el-text-color-regular);
}

.pa-hero__quote {
  margin-bottom: 20px;
}

.pa-hero__quote-from {
  display: block;
  margin-top: 6px;
  font-family: var(--pa-mono);
  font-size: 11px;
  letter-spacing: 0.1em;
  color: var(--el-text-color-secondary);
}

.pa-hero__note {
  font-size: 11px;
  line-height: 1.9;
  color: var(--el-text-color-secondary);
  max-width: 44em;
}

/* 叠放的曲绘 */

.pa-hero__art {
  position: relative;
  height: 300px;
}

.pa-hero__card {
  position: absolute;
  width: 62%;
  border: 1px solid var(--pa-line);
  border-radius: 5px;
  overflow: hidden;
  box-shadow: 0 18px 40px rgb(0 0 0 / 0.35);
  transition: transform 0.5s var(--gb-ease);
}

.pa-hero__card .pa-film {
  height: 150px;
}

.pa-hero__card-name {
  display: block;
  padding: 7px 10px;
  font-family: var(--pa-serif);
  font-size: 11px;
  letter-spacing: 0.06em;
  color: var(--el-text-color-secondary);
  background: var(--gb-card);
}

.pa-hero__card.is-0 {
  left: 0;
  top: 34px;
  transform: rotate(-6deg);
  z-index: 1;
}

.pa-hero__card.is-1 {
  left: 20%;
  top: 0;
  transform: rotate(2deg);
  z-index: 3;
}

.pa-hero__card.is-2 {
  left: 38%;
  top: 52px;
  transform: rotate(8deg);
  z-index: 2;
}

.pa-hero__art:hover .pa-hero__card.is-0 {
  transform: rotate(-11deg) translate(-10px, -4px);
}

.pa-hero__art:hover .pa-hero__card.is-2 {
  transform: rotate(13deg) translate(10px, 4px);
}

/* ---------- 时间轴 ---------- */

.pa-timeline {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 22px;
  border-top: 1px solid var(--pa-line);
  border-bottom: 1px solid var(--pa-line);
}

.pa-timeline__node {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.pa-timeline__node.is-center {
  text-align: center;
}

.pa-timeline__node.is-end {
  text-align: right;
}

.pa-timeline__date {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--gb-primary);
}

.pa-timeline__label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.pa-timeline__line {
  position: relative;
  flex: 1;
  height: 1px;
  background: var(--pa-line);
}

.pa-timeline__line i {
  position: absolute;
  left: 50%;
  top: 50%;
  width: 5px;
  height: 5px;
  margin: -2.5px 0 0 -2.5px;
  background: var(--gb-primary);
  transform: rotate(45deg);
}

/* ---------- 区块标题 ---------- */

.pa-section__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--pa-line);
  margin-bottom: 18px;
}

.pa-section__title {
  margin: 0;
  font-family: var(--pa-serif);
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--el-text-color-primary);
}

.pa-section__sub {
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
}

.pa-section__desc {
  margin: -6px 0 18px;
  font-size: 13px;
  line-height: 1.95;
  color: var(--el-text-color-regular);
  max-width: 62em;
}

/* ---------- 人物 ---------- */

.pa-cast {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.pa-cast__card {
  position: relative;
  padding: 18px;
  border: 1px solid var(--pa-line);
  border-left: 3px solid var(--pa-cast);
  border-radius: 5px;
  background: var(--gb-card);
  overflow: hidden;
}

.pa-cast__card::after {
  content: '';
  position: absolute;
  right: -30px;
  top: -30px;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  background: var(--pa-cast);
  opacity: 0.12;
  filter: blur(6px);
}

.pa-cast__head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.pa-cast__avatar {
  width: 34px;
  height: 34px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--pa-serif);
  font-size: 17px;
  font-weight: 700;
  color: #0b0a0e;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--pa-cast), var(--pa-cast-alt));
}

.pa-cast__name {
  font-family: var(--pa-serif);
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: var(--el-text-color-primary);
}

.pa-cast__vocal {
  font-size: 11px;
  letter-spacing: 0.12em;
  color: var(--el-text-color-secondary);
}

.pa-cast__role {
  display: inline-block;
  padding: 2px 8px;
  margin-bottom: 10px;
  font-size: 11px;
  letter-spacing: 0.08em;
  border-radius: 3px;
  color: var(--pa-cast);
  border: 1px solid color-mix(in srgb, var(--pa-cast) 40%, transparent);
  background: color-mix(in srgb, var(--pa-cast) 10%, transparent);
}

.pa-cast__line {
  margin: 0;
  font-size: 12px;
  line-height: 1.9;
  color: var(--el-text-color-secondary);
}

/* ---------- 筛选 ---------- */

.pa-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 18px;
}

.pa-filter {
  padding: 6px 14px;
  font-size: 12px;
  font-family: var(--pa-mono);
  letter-spacing: 0.06em;
  color: var(--el-text-color-secondary);
  border: 1px solid var(--pa-line);
  border-radius: 3px;
  background: transparent;
  cursor: pointer;
  transition:
    color 0.2s var(--gb-ease),
    border-color 0.2s var(--gb-ease),
    background 0.2s var(--gb-ease);
}

.pa-filter:hover {
  color: var(--gb-primary-deep);
  border-color: var(--pa-line-strong);
}

.pa-filter.is-active {
  color: #fff;
  border-color: var(--gb-primary);
  background: var(--gb-primary);
}

/* ---------- 卷宗卡 ---------- */

.pa-chapters {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(268px, 1fr));
  gap: 16px;
}

.pa-chapter {
  display: flex;
  flex-direction: column;
  padding: 0;
  border: 1px solid var(--pa-line);
  border-radius: 6px;
  background: var(--gb-card);
  text-align: left;
  cursor: pointer;
  overflow: hidden;
  transition:
    transform 0.35s var(--gb-ease),
    border-color 0.35s var(--gb-ease),
    box-shadow 0.35s var(--gb-ease);
}

.pa-chapter:hover {
  transform: translateY(-4px);
  border-color: var(--pa-line-strong);
  box-shadow: var(--gb-shadow-hover);
}

.pa-chapter__art {
  position: relative;
  aspect-ratio: 16 / 10;
}

.pa-chapter__sigil {
  position: absolute;
  left: 10px;
  top: 10px;
  z-index: 2;
  padding: 3px 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #fff;
  background: rgb(10 6 12 / 0.62);
  border: 1px solid rgb(255 255 255 / 0.22);
  border-radius: 3px;
  backdrop-filter: blur(4px);
}

.pa-chapter__body {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px 16px 16px;
}

.pa-chapter__meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  font-size: 11px;
  letter-spacing: 0.06em;
  color: var(--el-text-color-secondary);
}

.pa-chapter__title {
  margin: 2px 0 0;
  font-family: var(--pa-serif);
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.04em;
  color: var(--el-text-color-primary);
}

.pa-chapter__summary {
  margin: 0;
  font-size: 12px;
  line-height: 1.85;
  color: var(--el-text-color-regular);
}

.pa-chapter__lyric {
  margin: 0;
  font-family: var(--pa-serif);
  font-size: 12px;
  line-height: 1.7;
  color: var(--gb-primary-deep);
  opacity: 0.85;
}

.pa-chapter__more {
  margin-top: 4px;
  font-family: var(--pa-mono);
  font-size: 11px;
  letter-spacing: 0.1em;
  color: var(--el-text-color-secondary);
}

.pa-chapter:hover .pa-chapter__more {
  color: var(--gb-primary);
}

/* ---------- 病历卡 ---------- */

.pa-records {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 14px;
}

.pa-record {
  position: relative;
  padding: 16px 16px 18px;
  border: 1px dashed var(--pa-line);
  border-left: 2px solid var(--gb-primary);
  border-radius: 4px;
  background: var(--gb-card);
}

.pa-record__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.pa-record__index {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: var(--gb-primary);
}

.pa-record__chapter {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

.pa-record__term {
  margin: 0;
  font-family: var(--pa-serif);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: var(--el-text-color-primary);
}

.pa-record__alias {
  margin: 4px 0 10px;
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--el-text-color-placeholder);
}

.pa-record__note {
  margin: 0;
  font-size: 12px;
  line-height: 1.95;
  color: var(--el-text-color-regular);
}

/* ---------- 双栏 ---------- */

.pa-section--split {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 32px;
}

.pa-axioms {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.pa-axioms li {
  display: flex;
  gap: 12px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--pa-line-soft);
}

.pa-axioms li:last-child {
  border-bottom: none;
}

.pa-axioms span {
  flex: none;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--gb-primary);
}

.pa-axioms p {
  margin: 0;
  font-size: 13px;
  line-height: 1.9;
  color: var(--el-text-color-regular);
}

.pa-disc {
  display: flex;
  gap: 14px;
  padding: 12px 0;
  border-bottom: 1px solid var(--pa-line-soft);
}

.pa-disc:last-child {
  border-bottom: none;
}

.pa-disc__group {
  flex: none;
  width: 56px;
  font-family: var(--pa-serif);
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: var(--gb-primary);
}

.pa-disc__items {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
}

.pa-disc__item {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

/* ---------- 落款 ---------- */

.pa-colophon {
  padding: 26px 0 8px;
  text-align: center;
}

.pa-colophon__quote {
  margin: 20px 0 22px;
  font-family: var(--pa-serif);
  font-size: 15px;
  line-height: 2.2;
  letter-spacing: 0.04em;
  color: var(--el-text-color-primary);
}

.pa-colophon__staff {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px 16px;
  margin-bottom: 16px;
}

.pa-colophon__item {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

.pa-colophon__item em {
  font-style: normal;
  margin-right: 5px;
  color: var(--el-text-color-placeholder);
}

.pa-colophon__span {
  margin: 0 0 12px;
  font-family: var(--pa-mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  color: var(--gb-primary);
}

.pa-colophon__notice {
  margin: 0 auto;
  max-width: 60em;
  font-size: 11px;
  line-height: 1.9;
  color: var(--el-text-color-placeholder);
}

/* ---------- 卷宗弹窗 ---------- */

.pa-dossier__art {
  position: relative;
  aspect-ratio: 16 / 10;
  border: 1px solid var(--pa-line);
  border-radius: 5px;
  margin-bottom: 16px;
}

.pa-dossier__sigil {
  position: absolute;
  left: 12px;
  top: 12px;
  z-index: 2;
  padding: 4px 10px;
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  background: rgb(10 6 12 / 0.62);
  border: 1px solid rgb(255 255 255 / 0.22);
  border-radius: 3px;
}

.pa-dossier__meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px 14px;
  margin-bottom: 12px;
}

.pa-dossier__meta .muted {
  font-size: 12px;
}

.pa-dossier__voices {
  margin: 0 0 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.pa-dossier__detail {
  margin: 0 0 16px;
  font-size: 13px;
  line-height: 2;
  color: var(--el-text-color-regular);
}

.pa-dossier__lyric {
  white-space: pre-line;
  margin-bottom: 18px;
}

.pa-dossier__case {
  padding: 14px 16px;
  border: 1px dashed var(--pa-line);
  border-radius: 4px;
  background: var(--el-fill-color-lighter);
}

.pa-dossier__case-term {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 8px;
}

.pa-dossier__case-term span {
  font-size: 10px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
}

.pa-dossier__case-term strong {
  font-family: var(--pa-serif);
  font-size: 15px;
  letter-spacing: 0.06em;
  color: var(--gb-primary-deep);
}

.pa-dossier__case p {
  margin: 0;
  font-size: 12px;
  line-height: 1.95;
  color: var(--el-text-color-regular);
}

/* ---------- 响应式 ---------- */

@media (max-width: 1024px) {
  .pa-hero {
    grid-template-columns: 1fr;
  }

  .pa-hero__art {
    height: 250px;
    max-width: 520px;
  }
}

@media (max-width: 768px) {
  .pa-page {
    gap: 26px;
  }

  .pa-hero {
    padding: 20px 16px 24px;
  }

  .pa-hero__art {
    height: 210px;
  }

  .pa-hero__card .pa-film {
    height: 104px;
  }

  .pa-timeline {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    padding: 16px;
  }

  .pa-timeline__line {
    height: 12px;
    width: 1px;
    flex: none;
    margin-left: 3px;
  }

  .pa-timeline__line i {
    left: 50%;
    top: 50%;
  }

  /* 竖排时间轴里三条线的方向统一，避免左右错位 */
  .pa-timeline__node.is-center,
  .pa-timeline__node.is-end {
    text-align: left;
  }

  .pa-cast,
  .pa-chapters,
  .pa-records {
    grid-template-columns: 1fr;
  }

  .pa-colophon__quote {
    font-size: 14px;
  }
}
</style>
